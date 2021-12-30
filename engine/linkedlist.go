package engine

// this is a wrapper of linkedlist node for Batch
type nodeBatch struct {
	b    *Batch
	next *nodeBatch
}

// this is a linkedlist of Batch, which gonna be used as unbounded queue
//
// this implementation is NOT goroutine-safe
type llBatch struct {
	head *nodeBatch
	tail *nodeBatch
	next *nodeBatch
}

func (ll *llBatch) add(b *Batch) {
	nb := &nodeBatch{b: b}
	if ll.tail == nil {
		ll.head = nb
		ll.tail = nb
		return
	}
	ll.tail.next = nb
	ll.tail = nb
}

func (ll *llBatch) hasItem() bool {
	return ll.head != nil
}

func (ll *llBatch) get() *Batch {
	if ll.head == nil {
		return nil
	}
	b := ll.head.b
	if ll.head == ll.tail {
		ll.head = nil
		ll.tail = nil
	} else {
		ll.head = ll.head.next
	}
	return b
}
