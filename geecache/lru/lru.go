package lru

import "container/list"

type Value interface {
	Len() int
}

type Cache struct {
	maxBytes  int64                         // 最大长度
	nBytes    int64                         // 已使用的长度
	ll        *list.List                    // 双向链表
	cache     map[string]*list.Element      // 缓存字典
	OnEvicted func(key string, value Value) // 回掉函数
}

type entry struct {
	key   string
	value Value
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 获取元素并将其移动到队尾
func (c *Cache) Get(key string) (Value, bool) {
	if elem, ok := c.cache[key]; ok {
		c.ll.MoveToBack(elem)
		kv := elem.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// RemoveOldest 移除元素
func (c *Cache) RemoveOldest() {
	ele := c.ll.Front()
	if ele != nil {

		kv := ele.Value.(*entry)
		c.ll.Remove(ele)
		delete(c.cache, kv.key)
		ele = nil // 帮助GC进行回收

		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 增加元素
func (c *Cache) Add(key string, value Value) {
	if elem, ok := c.cache[key]; ok {
		et := elem.Value.(*entry)
		et.value = value
		c.nBytes += int64(value.Len()) - int64(elem.Value.(*entry).value.Len())
		c.ll.MoveToBack(elem)
	} else {
		c.cache[key] = c.ll.PushBack(&entry{key: key, value: value})
		c.nBytes += int64(len(key) + value.Len())
	}

	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
