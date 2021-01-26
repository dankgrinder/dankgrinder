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
	Client       *discord.Client
	Logger       *logrus.Logger
	ChannelID    string
	Typing       *config.Typing
	MessageDelay *config.MessageDelay

	priorityQueue *queue
	queue         *queue
	close         chan bool
	closed bool
	abort         chan bool
}

type Command struct {
	Run string

	// If not an empty string, this is what will be displayed by logrus when.
	// sending the command. The format will be "%v: %v", log, content.
	Log string

	// The interval at which the command should be rescheduled. Set to 0 to
	// disable.
	Interval time.Duration
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

	s.queue, s.priorityQueue = newQueue(), newQueue()
	s.close, s.abort = make(chan bool), make(chan bool, 1)

	s.priorityQueue.onEnqueue = func() {
		// Since the abort channel is buffered anyway, we do not want to remain
		// dormant attempting to send on abort. A scenario might also occur
		// where many aborts are sent around the same time, because of a busy
		// priority channel. This approach also prevent that.
		select {
		case s.abort <- true:
		default:
		}
	}

	go func() {
		for {
			// Clear the abort channel. Otherwise a scenario might occur
			// where an abort is sent during the execution of a priority
			// command and the next regular command is canceled due to that
			// abort, even though the priority command was already executed.
			select {
			case <-s.abort:
			case <-s.close:
				return
			default:
			}
			if s.priorityQueue.queued.Len() > 0 {
				cmd := <-s.priorityQueue.dequeue
				s.send(cmd, nil)
				continue
			}
			var cmd *Command
			select {
			case <-s.close:
				return
			case cmd = <-s.priorityQueue.dequeue:
				s.send(cmd, nil)
			case cmd = <-s.queue.dequeue:
				s.send(cmd, s.abort)
			}
		}
	}()
	return nil
}

func (s *Scheduler) Schedule(cmd *Command, priority bool) {
	if s.closed {
		return
	}
	if priority {
		s.priorityQueue.enqueue <- cmd
		return
	}
	s.queue.enqueue <- cmd
}

// Close closes the scheduler.
func (s *Scheduler) Close() error {
	s.closed = true
	s.close <- true
	close(s.close)
	if err := s.queue.Close(); err != nil {
		return err
	}
	if err := s.priorityQueue.Close(); err != nil {
		return err
	}
	close(s.abort)
	<-s.abort
	return nil
}

func (s *Scheduler) send(cmd *Command, abort chan bool) {
	d := delay(s.MessageDelay)
	tt := typing(cmd.Run, s.Typing)
	info := "sending command"
	if cmd.Log != "" {
		info = cmd.Log
	}
	s.Logger.WithFields(map[string]interface{}{
		"delay":  d.String(),
		"typing": tt.String(),
	}).Infof("%v: %v", info, cmd.Run)

	select {
	case <-abort:
		s.Logger.Infof("honoring abort of command: %v", cmd.Run)
		s.Schedule(cmd, false)
		return
	case <-time.After(d):
	}
	if err := s.Client.SendMessage(cmd.Run, discord.SendMessageOpts{
		ChannelID: s.ChannelID,
		Typing:    tt,
		Abort:     abort,
	}); err == discord.ErrAborted {
		s.Logger.Infof("honoring abort of command: %v", cmd.Run)
		s.Schedule(cmd, false)
		return
	} else if err != nil {
		s.Logger.Errorf("%v", err)
	}
	if cmd.Interval > 0 {
		time.AfterFunc(cmd.Interval, func() {
			s.Schedule(cmd, false)
		})
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
