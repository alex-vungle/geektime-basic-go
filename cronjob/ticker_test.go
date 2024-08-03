package cronjob

import (
	"context"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 间隔1秒的ticker
	ticker := time.NewTicker(time.Second)
	defer cancel()
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			t.Log("循环结束啦")
			goto end
		case now := <-ticker.C:
			t.Log("过了一秒", now.UnixMilli())
		}
	}
end:
	t.Log("来自goto，程序结束")
}
