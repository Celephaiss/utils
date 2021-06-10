package utils

import (
	"context"
	"log"
	"sync"
	"time"
)

type ApplyFunc = func(interface{}) interface{}

type ApplyFuncWithContext = func(context.Context, interface{}) interface{}

func Apply(collections []interface{}, f ApplyFunc) []interface{} {
	result := make([]interface{}, len(collections))
	wg := sync.WaitGroup{}
	wg.Add(len(collections))

	for i, v := range collections {
		go func(idx int, value interface{}) {
			result[idx] = f(value)
			wg.Done()
		}(i, v)
	}

	wg.Wait()
	return result
}

type entry struct {
	i   int
	arg interface{}
}

func Apply2(collections []interface{}, f ApplyFunc, nGoroutine int) []interface{} {

	if nGoroutine <= 0 {
		log.Panicf("num of goroutine less than 1")
	}

	result := make([]interface{}, len(collections))
	wg := sync.WaitGroup{}
	wg.Add(len(collections))

	ch := make(chan *entry)

	for i := 0; i < nGoroutine; i++ {
		go func() {
			for {
				entry, ok := <-ch
				if !ok {
					// the channel has been closed.
					break
				}
				result[entry.i] = f(entry.arg)
				wg.Done()

			}
		}()
	}

	for i, v := range collections {
		ch <- &entry{i, v}
	}
	close(ch)
	wg.Wait()
	return result
}

func Any(collections []interface{}, f ApplyFuncWithContext, nGoroutine int) (result interface{}) {

	if nGoroutine <= 0 {
		log.Panicf("num of goroutine less than 1")
	}

	resultCh := make(chan interface{})

	ch := make(chan *entry)

	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < nGoroutine; i++ {
		go func() {
			for {
				entry, ok := <-ch
				if !ok {
					// the channel has been closed.
					break
				}
				resultCh <- f(ctx, entry.arg)

			}
		}()
	}

	for i, v := range collections {
		ch <- &entry{i, v}
	}
	close(ch)
	result = <-resultCh
	cancel()
	return result
}

func All(collections []interface{}, f ApplyFuncWithContext, nGoroutine int, timeout time.Duration) (result []interface{}) {
	if nGoroutine <= 0 {
		log.Panicf("num of goroutine less than 1")
	}

	//resultCh := make(chan *entry)

	ch := make(chan *entry)

	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	wg.Add(len(collections))

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	for i := 0; i < nGoroutine; i++ {
		go func() {
			for {
				entry, ok := <-ch
				if !ok {
					// the channel has been closed.
					break
				}

				result[entry.i] = f(ctx, entry.arg)
				wg.Done()

			}
		}()
	}

	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	for i, v := range collections {
		ch <- &entry{i, v}
	}

	select {
	case <-time.After(timeout):
		cancel()
		return nil
	case <-doneCh:
		cancel()
		return result
	}
}
