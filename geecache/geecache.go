package geecache

import (
	"fmt"
	"sync"
)

var (
	rwmu   sync.RWMutex
	groups = make(map[string]*Group)
)

// Getter 回调函数接口
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 回调函数类型
type GetterFunc func(key string) ([]byte, error)

// Get 实现回调接口接口
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	mainCache cache
	getter    Getter
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	rwmu.Lock()
	defer rwmu.Unlock()

	g := &Group{
		name:      name,
		mainCache: cache{cacheBytes: cacheBytes},
		getter:    getter,
	}

	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	rwmu.RLock()
	defer rwmu.Unlock()
	return groups[name]
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("查询键不能为空")
	}

	value, ok := g.mainCache.get(key)
	if ok {
		return value, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
