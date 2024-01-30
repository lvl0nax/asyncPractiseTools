package batcher

import (
	"sync"
)

type Batcher struct {
	size     int
	action   func([]string)
	mutex    sync.Mutex
	messages []string

	batchesCh   chan []string
	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewBatcher(size int, action func([]string)) *Batcher {
	return &Batcher{
		size:        size,
		action:      action,
		batchesCh:   make(chan []string, 1),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}
}

func (b *Batcher) Append() {

}

func (b *Batcher) Run() {
	go func() {
		defer close(b.closeDoneCh)
		for {
			select {
			case <-b.closeCh:
				return
			case batch := <-b.batchesCh:
				b.action(batch)
			}
		}
	}()
}

func (b *Batcher) Shutdown() {
	close(b.closeCh)
	<-b.closeDoneCh
}
