package cronjob

import (
	cron "github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCronExpr(t *testing.T) {
	expr := cron.New(cron.WithSeconds())

	id, err := expr.AddFunc("@every 1s", func() {
		t.Log("执行了")
	})
	assert.NoError(t, err)
	t.Log("任务", id)
	expr.Start()
	time.Sleep(time.Second * 10)
	ctx := expr.Stop()
	t.Log("发出停止信号")
	<-ctx.Done()
	t.Log("没有任务执行了")

}

type JobFunc func()

func (j JobFunc) Run() {
	j()
}
