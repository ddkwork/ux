package safemap

import "container/list"

func (s *M[K, V]) Lock() {
	s.checkInit()
	s.lock.Lock()
}
func (s *M[K, V]) Unlock() { s.lock.Unlock() }
func (s *M[K, V]) RLock() {
	s.checkInit()
	s.lock.RLock()
}
func (s *M[K, V]) RUnlock() { s.lock.RUnlock() }

func (s *M[K, V]) checkInit() {
	if s.m == nil { // new(safemap.M[K, V])这种方式实例化代码简洁
		s.init(Ordered)
	}
}

func (s *M[K, V]) init(ordered ...bool) {
	*s = M[K, V]{
		m:        make(map[K]V),
		ordered:  false,
		keys:     list.New(),
		keyIndex: make(map[K]*list.Element),
	}
	s.keys.Init()
	if len(ordered) > 0 && ordered[0] {
		s.ordered = true
	}
}
