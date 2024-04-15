package sync

import (
	"sync"
	"testing"
)

// 假设说是两个 goroutine 交叉打印
func TestPrint100_2(t *testing.T) {
	var wg sync.WaitGroup
	// 两个任务
	wg.Add(2)
	ch := make(chan struct{})
	go func() {
		// 步长是 2
		for i := 0; i <= 100; i = i + 2 {
			<-ch
			t.Log("goroutine0", i)
			if i+1 <= 100 {
				// 发一个信号
				ch <- struct{}{}
			}
		}
		wg.Done()
	}()

	go func() {
		// 从 1 开始
		// 发一个信号开始
		ch <- struct{}{}
		for i := 1; i <= 100; i = i + 2 {
			<-ch
			t.Log("goroutine1", i)
			if i+1 <= 100 {
				ch <- struct{}{}
			}
		}
		wg.Done()
	}()
	wg.Wait()
}

func TestPrint100_3(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	ch0 := make(chan struct{})
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	go func() {
		for i := 0; i <= 100; i = i + 3 {
			<-ch0
			t.Log("goroutine0", i)
			if i+1 <= 100 {
				ch1 <- struct{}{}
			}
		}
		wg.Done()
	}()

	go func() {
		for i := 1; i <= 100; i = i + 3 {
			<-ch1
			t.Log("goroutine1", i)
			if i+1 <= 100 {
				ch2 <- struct{}{}
			}
		}
		wg.Done()
	}()

	go func() {
		// 发一个开始信号
		ch0 <- struct{}{}
		for i := 2; i <= 100; i = i + 3 {
			<-ch2
			t.Log("goroutine2", i)
			if i+1 <= 100 {
				ch0 <- struct{}{}
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
