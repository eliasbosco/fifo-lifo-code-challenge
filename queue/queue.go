package queue

import (
	unicorns "github.com/eliasbosco/fifo-lifo-code-challenge/unicorns/unicorns"
)

type QueueElement struct {
	Unicorns  *unicorns.UnicornList `json:"unicorns"`
	RequestId int                   `json:"request_id"`
	Status    *string               `json:"status"`
}

type QueueList []QueueElement

var Queue = &QueueList{}

func (q *QueueList) Enqueue(value *QueueElement) *QueueElement {
	queue := *q
	queue = append(queue, *value)
	*q = queue
	return &queue[len(queue)-1]
}

func (q *QueueList) Dequeue() *QueueElement {
	queue := *q
	if len(*q) > 0 {
		dequeued := queue[0]
		*q = queue[1:]
		return &dequeued
	}
	return nil
}

func (q *QueueList) FindByRequestIdFirstPosition(requestId int) *QueueElement {
	queue := *q
	for i, item := range queue {
		if i == 0 && item.RequestId == requestId {
			return &item
		} else {
			return nil
		}
	}
	return nil
}
