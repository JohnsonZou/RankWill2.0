package util

import "sync"

type Queue struct {
	items []interface{}
	mu    sync.Mutex
}

// 入队
func (q *Queue) Push(item interface{}) {
	q.mu.Lock()
	q.items = append(q.items, item)
	q.mu.Unlock()
}
func (q *Queue) Empty() bool {
	return len(q.items) == 0
}

// 出队
func (q *Queue) Pop() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}
