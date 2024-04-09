package sync

import "sync"

type Service interface {
	DoSomething()
}

// 单例模式：就是一个结构体，它只有一个实例
// 你这个东西，推荐用小写，你要是首字母大写，人家会绕开你的 GetInstance 方法自己创建实例
// &Singleton{}
type singleton struct {
}

func (s *singleton) DoSomething() {

}

var (
	instance *singleton
	initOnce = &sync.Once{}
)

// 这种是 lazy 模式
func GetInstance() Service {
	initOnce.Do(func() {
		instance = &singleton{}
	})
	return instance
}

// 还有一种饥渴模式
var (
	instance1 = &singleton{}
)

func GetInstance1() Service {
	return instance1
}

// 还有一种饥渴模式，初始化复杂结构体
var (
	instance2 *singleton
)

func init() {
	/// 一大堆代码
	instance2 = &singleton{}
}

func GetInstance2() Service {
	return instance2
}
