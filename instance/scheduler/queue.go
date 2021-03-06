// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package scheduler

import (
	"container/list"
)

type queue struct {
	enqueue chan *Command
	dequeue chan *Command

	// onEnqueue is a function executed when a new command is queued. Used
	// currently for sending an abort in the scheduler when a priority command
	// is enqueued.
	onEnqueue func()
	queued    *list.List
	close     chan struct{}
}

func newQueue() *queue {
	q := &queue{
		enqueue:   make(chan *Command),
		dequeue:   make(chan *Command),
		queued:    list.New(),
		onEnqueue: func() {},
		close:     make(chan struct{}),
	}
	go func() {
		for {
			if q.queued.Len() == 0 {
				select {
				case <-q.close:
					return
				case cmd := <-q.enqueue:
					q.queued.PushBack(cmd)
					go q.onEnqueue()
				}
				continue
			}
			select {
			case <-q.close:
				return
			case cmd := <-q.enqueue:
				q.queued.PushBack(cmd)
				go q.onEnqueue()
			case q.dequeue <- q.queued.Front().Value.(*Command):
				q.queued.Remove(q.queued.Front())
			}
		}
	}()
	return q
}

func (q *queue) Close() error {
	q.close <- struct{}{}
	close(q.close)

	// Make sure that all goroutines currently sending on q.enqueue don't panic
	// because it is closed. This assumes there will be no future senders.
	var done bool
	for !done {
		select {
		case <-q.enqueue:
		default:
			done = true
		}
	}
	close(q.enqueue)
	close(q.dequeue)
	return nil
}
