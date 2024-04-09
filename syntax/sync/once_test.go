package sync

import (
	"sync"
	"testing"
	"time"
)

func TestOnce(t *testing.T) {
	once := &sync.Once{}
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(func() {
				// 这句话会打印几次？
				// 只会打印一次
				t.Log("ABC")
			})
		}()
	}

	time.Sleep(time.Second)
}
