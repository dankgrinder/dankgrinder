// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package scheduler

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/dankgrinder/dankgrinder/config"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	Client             *discord.Client
	Logger             *logrus.Logger
	ChannelID          string
	Typing             *config.Typing
	MessageDelay       *config.MessageDelay
	AwaitResumeTimeout time.Duration
	FatalHandler       func(err error)

	queue              *queue
	priorityQueue      *queue
	close              chan struct{}
	closed             bool
	resume             chan *Command
	awaitResume        bool
	awaitResumeTrigger string
}

type Command struct {
	Value string

	// If not an empty string, this is what will be logged to the logger when
	// sending the command. The format will be "%v: %v", log, content.
	Log string

	// The interval at which the command should be rescheduled. Set to 0 to
	// disable.
	Interval time.Duration

	// If AwaitResume is true, the scheduler will wait for a resume call before
	// executing the next command.
	AwaitResume bool

	// Next is a pointer to the command that will be rescheduled if interval is
	// not 0. If this is nil this command itself will be rescheduled. Using this
	// feature might be useful when you want to create a chain of commands that
	// work together but have, for example, different values for Command.Value.
	Next *Command

	// The amount of times to reschedule the command in total. Set to 0 to
	// reschedule indefinitely. To run a command once, the interval should be
	// set to 0, not the amount.
	Amount uint

	execs uint
}

func (s *Scheduler) Start() error {
	if s.Client == nil {
		return fmt.Errorf("no client")
	}
	if s.ChannelID == "" {
		return fmt.Errorf("no channel id")
	}
	if s.Logger == nil {
		s.Logger = logrus.StandardLogger()
	}
	if s.AwaitResumeTimeout == 0 {
		s.AwaitResumeTimeout = math.MaxInt64
	}
	if s.FatalHandler == nil {
		s.FatalHandler = func(err error) {}
	}

	s.queue, s.priorityQueue = newQueue(), newQueue()
	s.close, s.resume = make(chan struct{}), make(chan *Command)

	go func() {
		for {
			if s.awaitResume {
				select {
				case cmd := <-s.resume:
					s.awaitResume = false
					if cmd != nil {
						s.send(cmd)
						continue
					}
				case <-time.After(s.AwaitResumeTimeout):
					s.awaitResume = false
				case <-s.close:
					return
				}
			}
			if s.priorityQueue.queued.Len() > 0 {
				cmd := <-s.priorityQueue.dequeue
				s.send(cmd)
				continue
			}
			select {
			case <-s.close:
				return
			case cmd := <-s.priorityQueue.dequeue:
				s.send(cmd)
			case cmd := <-s.queue.dequeue:
				s.send(cmd)
			}

		}
	}()
	return nil
}

// AwaitResumeTrigger returns the value of the command that caused the await
// resume state. An empty string will be returned if the scheduler is not awaiting
// a resume at the time this method is called.
func (s *Scheduler) AwaitResumeTrigger() string {
	if !s.awaitResume {
		return ""
	}
	return s.awaitResumeTrigger
}

func (s *Scheduler) Schedule(cmd *Command) {
	if s.closed {
		return
	}
	if cmd.Next == nil {
		cmd.Next = cmd
	}
	s.queue.enqueue <- cmd
}

func (s *Scheduler) PrioritySchedule(cmd *Command) {
	if s.closed {
		return
	}
	s.priorityQueue.enqueue <- cmd
}

// Resume makes a scheduler continue after being paused by a command with
// an AwaitResume value of true. Will block until scheduler has received the
// resume call.
func (s *Scheduler) Resume() {
	if s.closed {
		return
	}
	if !s.awaitResume {
		return
	}
	s.resume <- nil
}

// ResumeWithCommandOrPrioritySchedule is the same as ResumeWithCommand, but if
// the scheduler is not awaiting a resume, it will schedule the command in the
// priority queue instead.
func (s *Scheduler) ResumeWithCommandOrPrioritySchedule(cmd *Command) {
	if s.closed {
		return
	}
	if !s.awaitResume {
		s.PrioritySchedule(cmd)
	}
	s.resume <- cmd
}

// ResumeWithCommand is the same as Resume but executes the passed command
// immediately after resuming.
func (s *Scheduler) ResumeWithCommand(cmd *Command) {
	if s.closed {
		return
	}
	if !s.awaitResume {
		return
	}
	s.resume <- cmd
}

// Close closes the scheduler.
func (s *Scheduler) Close() error {
	s.closed = true
	s.close <- struct{}{}
	close(s.close)
	if err := s.queue.Close(); err != nil {
		return err
	}
	if err := s.priorityQueue.Close(); err != nil {
		return err
	}
	close(s.resume)
	return nil
}

// reschedule reschedules the command if the conditions for this are met. If so
// it will reschedule with the appropriate next command, based on the value of
// Command.Next.
func (s *Scheduler) reschedule(cmd *Command) {
	cmd.execs++
	if cmd.Amount != 0 && cmd.execs >= cmd.Amount {
		return
	}
	if cmd.Interval > 0 {
		time.AfterFunc(cmd.Interval, func() {
			if cmd.Next == nil {
				s.Schedule(cmd)
				return
			}
			s.Schedule(cmd.Next)
		})
	}
}

func (s *Scheduler) send(cmd *Command) {
	d := delay(s.MessageDelay)
	tt := typing(cmd.Value, s.Typing)
	info := "sending command"
	if cmd.Log != "" {
		info = cmd.Log
	}
	s.Logger.WithFields(map[string]interface{}{
		"delay":  d.String(),
		"typing": tt.String(),
	}).Infof("%v: %v", info, cmd.Value)

	if err := s.Client.SendMessage(cmd.Value, s.ChannelID, tt); err == discord.ErrForbidden || err == discord.ErrInvalidAuthorization {
		s.FatalHandler(err)
		// Ran in a goroutine because otherwise the scheduler's goroutine would
		// be attempting to send a message to itself via its close channel which
		// just causes a permanently dormant goroutine.
		go s.Close()

		// Set to true to make sure the scheduler doesn't loop back around so
		// fast that it hasn't been closed yet by the goroutine created previously.
		s.awaitResume = true
		return
	} else if err != nil {
		s.Logger.Errorf("%v", err)
		s.Schedule(cmd)
		return
	}
	s.reschedule(cmd)
	if cmd.AwaitResume {
		s.awaitResumeTrigger = cmd.Value
		s.awaitResume = true
	}
}

// typing returns a duration for which to type based on the variables in the
// config.
func typing(cmd string, typing *config.Typing) time.Duration {
	msPerKey := int(math.Round((1.0 / float64(typing.Speed)) * 60000))
	d := typing.Base
	d += len(cmd) * msPerKey
	if typing.Variance > 0 {
		d += rand.Intn(typing.Variance)
	}
	return time.Duration(d) * time.Millisecond
}

// delay returns a duration for which to sleep before commencing typing based on
// the variables in the config.
func delay(messageDelay *config.MessageDelay) time.Duration {
	d := messageDelay.Base
	if messageDelay.Variance > 0 {
		d += rand.Intn(messageDelay.Variance)
	}
	return time.Duration(d) * time.Millisecond
}
