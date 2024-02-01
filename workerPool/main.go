package workerPool

import (
	"errors"
	"sync"
)

type WorkerPool struct {
	tasksCh     chan func()
	closeCh     chan struct{}
	closeDoneCh chan struct{}

	mutex sync.RWMutex
}

func NewWorkerPool(size int) *WorkerPool {
	wp := &WorkerPool{
		tasksCh:     make(chan func(), size),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}

	wp.initializeWorkers(size)
	return wp
}

func (wp *WorkerPool) AddTask(task func()) error {
	if task == nil {
		return errors.New("task is nil")
	}

	wp.mutex.RLock()
	defer wp.mutex.RUnlock()

	select {
	case <-wp.closeCh:
		return errors.New("worker pool is closed")
	default:
	}

	select {
	case wp.tasksCh <- task:
		return nil
	default:
		return errors.New("worker pool is full")
	}
}

func (wp *WorkerPool) ShutDown() {
	close(wp.closeCh)

	wp.mutex.Lock()
	close(wp.tasksCh)
	wp.mutex.Unlock()

	<-wp.closeDoneCh
}

func (wp *WorkerPool) initializeWorkers(size int) {
	wg := sync.WaitGroup{}
	wg.Add(size)

	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			for task := range wp.tasksCh {
				task()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(wp.closeDoneCh)
	}()
}
