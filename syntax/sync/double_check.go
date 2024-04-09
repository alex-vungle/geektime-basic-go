package sync

import "sync"

// 你们也可以做成泛型
type SafeMap struct {
	data map[string]string
	lock sync.RWMutex
}

// LoadOrStore 如果 key 存在，就返回老的数据，并且 loaded 返回 true
// 否则，把 key 的值设置为 newVal，返回 newVal，并且 loaded 返回 false
func (m *SafeMap) LoadOrStore(key string, newVal string) (val string, loaded bool) {
	// 最粗暴的做法
	//m.lock.Lock()
	//defer m.lock.Unlock()

	// double check 写法
	m.lock.RLock()
	oldVal, ok := m.data[key]
	if ok { // 读多写少，也就是大部分请求命中这个分支，就用 double -check 写法
		m.lock.RUnlock()
		return oldVal, true
	}
	m.lock.RUnlock()
	// 这样写有没有问题？
	m.lock.Lock()
	defer m.lock.Unlock()

	// 这个就是最关键的，第二次检查，
	// double check 就是指检查两次
	oldVal, ok = m.data[key]
	if ok {
		return oldVal, true
	}

	m.data[key] = newVal
	return newVal, false
	// 总结，所有的 检查-做某事 类的问题，你就是两种方案
	// 1. 直接写锁，检查，做某事。适合写多读少
	// 2. 先读锁，检查，符合条件返回；不符合条件，加写锁，再检查，做某事。适合读多写少

	// 分布式环境下的检查-做某事也是一样的
}
