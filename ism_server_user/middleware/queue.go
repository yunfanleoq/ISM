/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:25
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:58:20
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.

The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/
package middleware

import "sync"

// minQueueLen is smallest capacity that queue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minQueueLen = 10000 // 队列缓存区最小长度

// Queue represents a single instance of the queue data structure.
type Queue struct {
	head int
	foot int
	arr  []interface{}
	lock *sync.Mutex
	cap  int
}

// New constructs and returns a new Queue.
func NewQueue(Length int) *Queue {
	if Length == 0 {
		Length = minQueueLen
	}
	return &Queue{
		head: 0, // 包含此下表值
		foot: 0, // 不包含此下标值
		arr:  make([]interface{}, Length+1),
		lock: &sync.Mutex{},
		cap:  Length,
	}
}

// Length returns the number of elements currently stored in the queue.
func (q *Queue) QueueLength() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return (q.foot + len(q.arr) - q.head) % len(q.arr)
}

// Add puts an element on the end of the queue.
func (q *Queue) QueuePush(elem interface{}) int {

	q.lock.Lock()
	defer q.lock.Unlock()
	if (q.foot+1)%len(q.arr) == q.head {
		return -1
	}
	q.arr[q.foot] = elem
	q.foot = (q.foot + 1) % len(q.arr)
	return 0
}

// Remove removes and returns the element from the front of the queue. If the
// queue is empty, the call will panic.
func (q *Queue) QueuePull() (interface{}, int) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.foot == q.head {
		return nil, -1
	}
	val := q.arr[q.head]
	q.head = (q.head + 1) % len(q.arr)
	return val, 0
}
