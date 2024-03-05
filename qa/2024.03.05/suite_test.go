package qa

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type MySuite struct {
	suite.Suite
}

func (m *MySuite) SetupSuite() {
	// 启动的时候执行
}

func (m *MySuite) TearDownSuite() {
	// 停下的时候执行
}

func (m *MySuite) SetupTest() {
	// 每一个测试开始的时候执行
}
func (m *MySuite) TearDownTest() {
	// 每一个测试执行完毕的时候执行
}

// 具体测试
func (m *MySuite) TestXXX() {

}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MySuite))
}
