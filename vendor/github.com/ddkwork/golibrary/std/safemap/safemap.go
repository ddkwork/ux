package safemap

import (
	"container/list"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"sync"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/mylog/pretty"
)

const Ordered = true

type M[K comparable, V any] struct {
	lock     sync.RWMutex
	m        map[K]V
	ordered  bool
	keys     *list.List
	keyIndex map[K]*list.Element
}

// thx github.com/hitsumitomo/safemap
// type api[K comparable, V any] interface {
//	NewContent(ordered ...bool) (m *M[K, V])                                      // 实例化
//	NewOrdered(yield iter.Seq2[K, V]) (m *M[K, V])                         // 实例化有序，std map的代码有语法检查，这个是实例化的时候检查，实例化语法间接性差不多
//	NewStringer(ordered ...bool) (m *M[string, string])                    // 从字符串实例化
//	NewStringerKeys(keys []string, ordered ...bool) (m *M[string, string]) // 从字符串切片实例化
//	Has(key K) (exists bool)                                               // 是否存在
//	Get(key K) (value V, exist bool)                                       // 获取
//	Update(key K, value V)                                                 // 更新
//	Delete(key K)                                                          // 删除
//	Remove(key K)                                                          // 移除key
//	Set(key K, value V) (actual V, exist bool)                             // 设置，如果存在则不更新
//	GetAndDelete(key K) (value V, exist bool)                              // 获取后删除
//	removeKey(key K)                                                       // 移除key
//	Range() iter.Seq2[K, V]                                                // 遍历，todo 回调内执行删除会死锁,vt调试器bind的时候需要
//	Reset()                                                                // 清空
//	Len() int                                                              // 大小
//	Empty() bool                                                           // 大小为0
//	Keys() []K                                                             // 键列表
//	Map() map[K]V                                                          // 原始map
//	CopyFromMap(data map[K]V)                                              // 从map复制
//	MarshalJSON() (data []byte, err error)                                 //
//	UnmarshalJSON(data []byte) (err error)                                 //
//	String() string                                                        //
// }

func New[K comparable, V any](ordered ...bool) (m *M[K, V]) {
	sm := &M[K, V]{
		m:        make(map[K]V),
		keys:     list.New(),
		keyIndex: make(map[K]*list.Element),
	}
	if len(ordered) > 0 && ordered[0] {
		sm.ordered = true
	}
	return sm
}

func (s *M[K, V]) Collect(seq iter.Seq2[K, V]) *M[K, V] {
	s.Lock()
	defer s.Unlock()
	for k, v := range seq {
		_, exist := s.Set(k, v)
		if exist {
			panic("duplicate key: " + fmt.Sprint(k))
		}
	}
	return s
}

func (s *M[K, V]) Reset() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[K]V)
	s.keys.Init()
	s.keyIndex = make(map[K]*list.Element)
}

func NewOrdered[K comparable, V any](yield iter.Seq2[K, V]) (m *M[K, V]) {
	m = New[K, V](true)
	for k, v := range yield {
		_, exist := m.Set(k, v)
		if exist {
			panic("duplicate key: " + fmt.Sprint(k))
		}
	}
	return
}

func NewStringer(ordered ...bool) (m *M[string, string]) {
	return New[string, string](ordered...)
}

func NewStringerKeys(keys []string, ordered ...bool) (m *M[string, string]) {
	m = NewStringer(ordered...)
	for _, key := range keys {
		m.Set(key, key)
	}
	return
}

func (s *M[K, V]) Has(key K) (exists bool) {
	s.RLock()
	defer s.RUnlock()
	_, exists = s.m[key]
	return exists
}

func (s *M[K, V]) GetMust(key K) (value V) {
	get, exist := s.Get(key)
	if !exist {
		mylog.Check("key: " + fmt.Sprint(key) + " not found")
	}
	return get
}

func (s *M[K, V]) GetMustCallback(key K, callback func(value V) V) (value V) {
	return callback(s.GetMust(key))
}

func (s *M[K, V]) Get(key K) (value V, exist bool) {
	s.RLock()
	defer s.RUnlock()
	value, exist = s.m[key]
	return value, exist
}

func (s *M[K, V]) Update(key K, value V) {
	s.Lock()
	defer s.Unlock()
	if elem, exists := s.keyIndex[key]; !exists {
		s.keyIndex[key] = s.keys.PushBack(key)
	} else {
		elem.Value = key
	}
	s.m[key] = value
}

func (s *M[K, V]) Set(key K, value V) (actual V, exist bool) {
	s.Lock()
	defer s.Unlock()
	actual, exist = s.m[key]
	if exist {
		mylog.CheckIgnore("map set exist key : " + fmt.Sprint(key))
		return actual, true
	}
	s.m[key] = value
	if s.ordered {
		if elem, exists := s.keyIndex[key]; !exists {
			s.keyIndex[key] = s.keys.PushBack(key)
		} else {
			elem.Value = key
		}
	}
	return value, false
}

func (s *M[K, V]) Range() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if s.ordered {
			for e := s.keys.Front(); e != nil; e = e.Next() {
				k := e.Value.(K)
				if !yield(k, s.m[k]) {
					return
				}
			}
			return
		}
		for k, v := range maps.All(s.m) {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (s *M[K, V]) RangeKeys() iter.Seq[K] { return maps.Keys(s.m) } // todo test 这种不支持有序遍历keys，使用Keys() 老方案遍历算了

func (s *M[K, V]) Keys() []K {
	s.RLock()
	defer s.RUnlock()
	keys := make([]K, s.keys.Len())
	i := 0
	if s.ordered {
		for e := s.keys.Front(); e != nil; e = e.Next() {
			keys[i] = e.Value.(K)
			i++
		}
	} else {
		for k := range maps.Keys(s.m) {
			keys[i] = k
		}
	}
	return keys
}

func (s *M[K, V]) Values() []V {
	s.RLock()
	defer s.RUnlock()
	values := make([]V, s.keys.Len())
	i := 0
	if s.ordered {
		for e := s.keys.Front(); e != nil; e = e.Next() {
			values[i] = s.m[e.Value.(K)]
			i++
		}
	} else {
		for v := range maps.Values(s.m) {
			values[i] = v
		}
	}
	return values
}

func (s *M[K, V]) Empty() bool  { return s.Len() == 0 }
func (s *M[K, V]) Remove(key K) { s.removeKey(key) }
func (s *M[K, V]) Delete(key K) { s.removeKey(key) }
func (s *M[K, V]) removeKey(key K) {
	s.Lock()
	defer s.Unlock()
	if elem, exists := s.keyIndex[key]; exists {
		s.keys.Remove(elem)
		delete(s.keyIndex, key)
	}
	delete(s.m, key)
}

func (s *M[K, V]) GetAndDelete(key K) (value V, exist bool) {
	value, exist = s.m[key]
	if exist {
		s.removeKey(key)
	}
	return value, exist
}

func (s *M[K, V]) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

func (s *M[K, V]) Map() map[K]V {
	s.RLock()
	defer s.RUnlock()
	return maps.Clone(s.m)
}

func (s *M[K, V]) CopyFromMap(data map[K]V) {
	s.Lock()
	defer s.Unlock()
	s.m = maps.Clone(data) // todo add Clone method,use deepcopy pkg
	s.keys.Init()
	s.keyIndex = make(map[K]*list.Element)
	for k := range data {
		s.keyIndex[k] = s.keys.PushBack(k)
	}
}

func (s *M[K, V]) MarshalJSON() (data []byte, err error) {
	s.RLock()
	defer s.RUnlock()
	return json.Marshal(s.m)
}

func (s *M[K, V]) UnmarshalJSON(data []byte) (err error) {
	s.Lock()
	defer s.Unlock()
	var m map[K]V
	mylog.Check(json.Unmarshal(data, &m))
	s.m = m
	s.keys.Init()
	s.keyIndex = make(map[K]*list.Element)
	return nil
}

func (s *M[K, V]) String() string {
	s.RLock()
	defer s.RUnlock()
	return pretty.Format(s.m)
}

func (s *M[K, V]) LastKey() K { return s.Keys()[s.Len()-1] }
func (s *M[K, V]) LastValue() V {
	return s.GetMust(s.LastKey())
}
