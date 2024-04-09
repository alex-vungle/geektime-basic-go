package sync

import "sync"

// Queue 你们可以改泛型
type Queue struct {
	lock     *sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	data     []any
	capacity int
}

func NewQueue(capacity int) *Queue {
	lock := &sync.Mutex{}
	return &Queue{
		data:     make([]any, 0, capacity),
		lock:     lock,
		notEmpty: sync.NewCond(lock),
		notFull:  sync.NewCond(lock),
	}
}

// Enqueue 有一个基本的逻辑，就是如果你已经满了，就阻塞等待
func (q *Queue) Enqueue(data any) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for q.capacity == len(q.data) {
		// 你已经满了，怎么阻塞住呢？

		// 我在这里等一个不满的信号
		// Wait 这里已经释放了锁
		q.notFull.Wait()

		// 你刚睡醒，结果别人把你的位置抢了
	}

	// 当你到这里的时候，你可以断定，这个位置被你抢到了

	q.data = append(q.data, data)
	// 唤醒一个等数据的人 Dequeue
	q.notEmpty.Signal()
}

// Dequeue 如果你已经空了，就阻塞等待
func (q *Queue) Dequeue() any {
	q.lock.Lock()
	defer q.lock.Unlock()
	for len(q.data) == 0 {
		// 空的，没有元素
		// 你得等
		q.notEmpty.Wait()
	}
	val := q.data[0]
	q.data = q.data[1:]
	// 你已经取走了一个元素，告诉对面的被阻塞的（可能有）
	q.notFull.Signal()
	return val
}
