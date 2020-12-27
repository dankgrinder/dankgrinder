package main

import (
	"container/list"
	"github.com/dankgrinder/dankgrinder/discord"
	"github.com/sirupsen/logrus"
	"time"
)

type scheduler struct {
	schedule      chan command
	priority      chan command
	priorityQueue queue
	queue         queue
}

type queue struct {
	enqueue chan command
	dequeue chan command
	queued  *list.List
}

type command struct {
	content string

	// The interval at which the command should be rescheduled. Set to 0 to
	// disable.
	interval time.Duration
}

func sendMessage(content string, abort chan bool) {
	d := delay()
	tt := typingTime(content)
	logrus.WithFields(map[string]interface{}{
		"delay":  d.String(),
		"typing": tt.String(),
	}).Infof("sending command: %v", content)
	time.Sleep(d)

	if err := auth.SendMessage(content, discord.SendMessageOpts{
		ChannelID:  cfg.ChannelID,
		TypingTime: tt,
		Abort:      abort,
	}); err != nil {
		logrus.Errorf("%v", err)
	}
}

func startNewQueue() queue {
	q := queue{
		enqueue: make(chan command),
		dequeue: make(chan command),
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
			case q.dequeue <- q.queued.Front().Value.(command):
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
				sendMessage(cmd.content, nil)
				if cmd.interval > 0 {
					s.reschedule(cmd)
				}
				continue
			}
			select {
			case cmd := <-s.priorityQueue.dequeue:
				sendMessage(cmd.content, nil)
				if cmd.interval > 0 {
					s.reschedule(cmd)
				}

			case cmd := <-s.queue.dequeue:
				sendMessage(cmd.content, abort)
				if cmd.interval > 0 {
					s.reschedule(cmd)
				}
			}
		}
	}()
	return s
}

func (s scheduler) reschedule(cmd command) {
	go func() {
		time.Sleep(cmd.interval)
		s.schedule <- cmd
	}()
}
