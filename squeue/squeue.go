package squeue

import (
	"runtime"
	"sync/atomic"
	"time"
)

const (
	lock = 1
	free = 0
)

type sNode struct {
	data interface{}
	next *sNode
}

// newNode 创建一个新的节点
func newNode(Data interface{}) *sNode {
	return &sNode{
		data: Data,
		next: nil,
	}
}

// SQueue 是一个高效
type SQueue struct {
	tail *sNode
	head *sNode
	lock int32
	len  uint64
}

// NewSDQueue 用于创建一个SD队列
func NewSQueue() *SQueue {
	PNode := newNode(nil)
	return &SQueue{
		head: PNode,
		tail: PNode,
		lock: free,
		len:  0,
	}
}

// isEmpty 判断这个队列是不是空的

func (sq *SQueue) isEmpty() bool {
	return sq.head == sq.tail
}

// Push 用于向队列追加元素
func (sq *SQueue) Push(Data interface{}) {
	PNewNode := newNode(Data)
	for {
		OldLock := atomic.LoadInt32(&sq.lock)
		if OldLock == lock {
			// 如果满足,则代表此时该队列是一个上了锁的队列
			runtime.Gosched()
			time.Sleep(time.Microsecond)
			continue
		}
		// 如果进入了这里,那么此时代表OldLock 是处于自由状态
		// 所以下面要对其进行上锁
		if !atomic.CompareAndSwapInt32(&sq.lock, OldLock, lock) {
			// 如果进入了这里,代表着上锁失败
			runtime.Gosched()
			time.Sleep(time.Microsecond)
			continue
		}
		// 进入到了这里代表着上锁成功了,此时可以添加元素
		sq.tail.next = PNewNode
		sq.tail = PNewNode
		sq.len++
		// 此时应该进行解锁
		atomic.CompareAndSwapInt32(&sq.lock, lock, free)
		break
	}
}

// Length 用于返回该队列的长度
func (sq *SQueue) Length() uint64 {
	return sq.len
}

// Pop 用于删除队列的首元素,并且返回该元素的数据
func (sq *SQueue) Pop() (Data interface{}) {
	for {
		OldLock := atomic.LoadInt32(&sq.lock)
		if OldLock == lock {
			// 如果满足,则代表此时该队列是一个上了锁的队列
			runtime.Gosched()
			time.Sleep(time.Microsecond)
			continue
		}
		// 如果进入了这里,那么此时代表OldLock 是处于自由状态
		// 所以下面要对其进行上锁
		if !atomic.CompareAndSwapInt32(&sq.lock, OldLock, lock) {
			// 如果进入了这里,代表着上锁失败
			runtime.Gosched()
			time.Sleep(time.Microsecond)
			continue
		}
		First := sq.head.next
		if First == nil {
			Data = nil
		} else {
			Data = First.data
			sq.head.next = First.next
			sq.len--
			if sq.head.next == nil {
				sq.tail = sq.head
				sq.head.next = sq.tail
			}
		}
		atomic.CompareAndSwapInt32(&sq.lock, lock, free)
		break
	}
	return Data
}
