package qa

import (
	"github.com/ecodeclub/ekit/slice"
	"testing"
)

func TestAbc(t *testing.T) {
	val := []int64{1, 2, 3}
	// 删除
	res := slice.FilterDelete(val, func(idx int, src int64) bool {
		return src%2 == 1
	})
	t.Log(res)
}
