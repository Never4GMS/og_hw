package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskChan := make(chan Task)
	wg := sync.WaitGroup{}
	mMutex := sync.Mutex{}
	checkErrors := m > 0
	errorLimitExceeded := notCheckingErrors

	if checkErrors {
		errorLimitExceeded = createErrorLimitExceeded(&m, &mMutex)
	}

	wg.Add(n + 1)
	go func(out chan<- Task) {
		defer wg.Done()
		defer close(out)
		for _, task := range tasks {
			if errorLimitExceeded() {
				break
			}
			out <- task
		}
	}(taskChan)

	for i := 0; i < n; i++ {
		go func(tasks <-chan Task) {
			defer wg.Done()
			for task := range tasks {
				err := task()
				if err != nil {
					mMutex.Lock()
					m--
					mMutex.Unlock()
				}
			}
		}(taskChan)
	}

	wg.Wait()

	if checkErrors && m <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func createErrorLimitExceeded(m *int, mMutex *sync.Mutex) func() bool {
	return func() (exceeded bool) {
		mMutex.Lock()
		exceeded = *m <= 0
		mMutex.Unlock()
		return exceeded
	}
}

func notCheckingErrors() bool {
	return false
}
