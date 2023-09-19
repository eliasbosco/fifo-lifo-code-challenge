package queue

import (
	"testing"
	queue "unicorn/queue"
)

var Queue = &queue.QueueList{}

func TestEnqueue(t *testing.T) {
	_queue := new(queue.QueueElement)
	_queue = Queue.Enqueue(_queue)
	_queue.RequestId = 1
	_queue = new(queue.QueueElement)
	_queue = Queue.Enqueue(_queue)
	_queue.RequestId = 2

	__queue := *Queue
	if __queue[0].RequestId != 1 {
		t.Fatalf(`Queue.Enqueue first position has no value RequestId: 1`)
	}
	if __queue[1].RequestId != 2 {
		t.Fatalf(`Queue.Enqueue first position has no value RequestId: 2`)
	}
}

func TestDequeue(t *testing.T) {
	_queue := new(queue.QueueElement)
	_queue = Queue.Enqueue(_queue)
	_queue.RequestId = 1
	_queue = new(queue.QueueElement)
	_queue = Queue.Enqueue(_queue)
	_queue.RequestId = 2

	__queue := *Queue

	if __queue[0].RequestId != 1 {
		t.Fatalf(`Queue.Enqueue first position has no value RequestId: 1`)
	}

	if __queue[1].RequestId != 2 {
		t.Fatalf(`Queue.Enqueue first position has no value RequestId: 2`)
	}

	// FIFO principle, first position dequeued
	_dequeued := *Queue.Dequeue()
	if _dequeued.RequestId != 1 {
		t.Fatalf(`Queue.Dequeue first position has no value RequestId: 1`)
	}

	// Old second position becomes first
	newQueue := *Queue
	if newQueue[0].RequestId != 2 {
		t.Fatalf(`Queue.Dequeue first position has no value RequestId: 2`)
	}
}
