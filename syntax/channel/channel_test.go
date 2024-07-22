package channel

import (
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	// 声明
	// var ch chan struct{}
	// 声明并初始化
	//ch1 := make(chan int)
	// 这种是带buffer的
	ch2 := make(chan int, 3)
	ch2 <- 123
	ch2 <- 456
	ch2 <- 789
	t.Log(<-ch2)
	close(ch2)
	t.Log(<-ch2)
	ch2 <- 988
}

func TestChannelForLoop(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
			time.Sleep(time.Second)
		}
		close(ch)
	}()
	for val := range ch {
		t.Log(val)
	}
	t.Log("发送完毕")
}

func TestChannelSelect(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 2)
	go func() {
		time.Sleep(time.Second)
		ch1 <- 123
	}()

	go func() {
		time.Sleep(time.Second)
		ch2 <- 234
	}()

	select {
	case val := <-ch1:
		t.Log("进来了ch1这里", val)
	case val := <-ch2:
		t.Log("进来了ch2这里", val)
	}
}
