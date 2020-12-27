package main

import (
	"container/list"
	"time"
)

type scheduler struct {
	schedule chan string
	priority chan string
	priorityQueue queue
	queue    queue
}

type queue struct {
	enqueue      chan string
	dequeue      chan string
	queued       *list.List
}

func startNewQueue() queue {
	q := queue{
		enqueue:      make(chan string),
		dequeue:      make(chan string),
		queued:       list.New(),
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
			case q.dequeue <- q.queued.Front().Value.(string):
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
		priority: qp.enqueue,
		schedule: q.enqueue,
		queue:    q,
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
				sendMessage(cmd, nil)
				continue
			}
			select {
			case cmd := <-s.priorityQueue.dequeue: sendMessage(cmd, abort)
			case cmd := <-s.queue.dequeue: sendMessage(cmd, abort)
			}
		}
	}()
	return s
}

func (s scheduler) scheduleInterval(cmd string, interval time.Duration) {
	t := time.Tick(interval)
	go func() {
		for {
			s.schedule <- cmd
			<-t
		}
	}()
}
