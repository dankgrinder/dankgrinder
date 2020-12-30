// Copyright (C) 2020 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"container/list"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/sirupsen/logrus"
	"time"
)

type scheduler struct {
	schedule      chan *command
	priority      chan *command
	priorityQueue queue
	queue         queue
}

type queue struct {
	enqueue chan *command
	dequeue chan *command
	queued  *list.List
}

type command struct {
	content string

	// If not an empty string, this is what will be displayed by logrus when.
	// sending the command. The format will be "%v: %v", log, content.
	log string

	// The interval at which the command should be rescheduled. Set to 0 to
	// disable.
	interval time.Duration
}

func startNewQueue() queue {
	q := queue{
		enqueue: make(chan *command),
		dequeue: make(chan *command),
		queued:  list.New(),
	}
	go func() {
		for {
			if q.queued.Len() == 0 {
				cmd := <-q.enqueue
				q.queued.PushBack(cmd)
				continue
			}
			select {
			case cmd := <-q.enqueue:
				q.queued.PushBack(cmd)
			case q.dequeue <- q.queued.Front().Value.(*command):
				q.queued.Remove(q.queued.Front())
			}
		}
	}()
	return q
}

func startNewScheduler() scheduler {
	q := startNewQueue()
	qp := startNewQueue()
	s := scheduler{
		priority:      qp.enqueue,
		schedule:      q.enqueue,
		queue:         q,
		priorityQueue: qp,
	}

	// An abort will be sent to free up the scheduler for a priority command.
	abort := make(chan bool)
	go func() {
		for {
			if s.priorityQueue.queued.Len() > 0 {
				abort <- true
			}
			time.Sleep(time.Millisecond)
		}
	}()

	go func() {
		for {
			if s.priorityQueue.queued.Len() > 0 {
				cmd := <-s.priorityQueue.dequeue
				s.send(cmd, nil)

				// Clear the abort channel. Otherwise a scenario might occur
				// where an abort is sent during the execution of a priority
				// command and the next regular command is canceled due to that
				// abort, even though the priority command was already executed.
				select {
				case <-abort:
				default:
				}
				continue
			}
			var cmd *command
			select {
			case cmd = <-s.priorityQueue.dequeue:
				s.send(cmd, nil)
			case cmd = <-s.queue.dequeue:
				s.send(cmd, abort)
			}
		}
	}()
	return s
}

func (s scheduler) send(cmd *command, abort chan bool) {
	d := delay()
	tt := typing(cmd.content)
	info := "sending command"
	if cmd.log != "" {
		info = cmd.log
	}
	logrus.WithFields(map[string]interface{}{
		"delay":  d.String(),
		"typing": tt.String(),
	}).Infof("%v: %v", info, cmd.content)

	select {
	case <-abort:
		return
	case <-time.After(d):
	}
	if err := auth.SendMessage(cmd.content, discord.SendMessageOpts{
		ChannelID: cfg.ChannelID,
		Typing:    tt,
		Abort:     abort,
	}); err == discord.ErrAborted {
		s.schedule <- cmd
		return
	} else if err != nil {
		logrus.Errorf("%v", err)
	}
	if cmd.interval > 0 {
		s.reschedule(cmd)
	}
}

// reschedule is run after a command has been sent by the scheduler to
// add the command to the back of the queue again. This should only be run
// by the scheduler!
func (s scheduler) reschedule(cmd *command) {
	go func() {
		time.Sleep(cmd.interval)
		s.schedule <- cmd
	}()
}
