package queue

import (
	unicorns "unicorn/unicorns"
)

const (
	PRODUCTION_STATUS_IN_PROGRESS string = "in progress"
	PRODUCTION_STATUS_READY              = "ready"
)

type QueueElement struct {
	Unicorns  *unicorns.UnicornList `json:"unicorns"`
	RequestId int                   `json:"request_id"`
	Status    *string               `json:"status"`
}

type QueueList []QueueElement

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

func (q *QueueList) FindQueueFirstPosition(requestId int) *QueueElement {
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

func (q *QueueList) FindQueueByRequestId(requestId int) *QueueElement {
	for _, item := range *q {
		if item.RequestId == requestId {
			return &item
		}
	}
	return nil
}
