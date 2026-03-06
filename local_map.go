package cache

import (
	"fmt"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/monopolly/cast"
)

type localMap struct {
	cache map[any]any
	ttl   map[any]int64
	m     sync.RWMutex
	max   int
	add   time.Duration
	name  string
}

// simple map
func NewMap(name string, ttl time.Duration) (res Engine) {
	a := new(localMap)
	a.cache = map[any]any{}
	a.ttl = map[any]int64{}
	a.add = ttl
	a.name = name
	go a.deamon()
	return a
}

func (a *localMap) sid(id any) (res string) {
	return sid(id, a.name) //users:id
}

func (a *localMap) set(id any, v any) {
	a.m.Lock()
	a.cache[a.sid(id)] = v
	a.ttl[a.sid(id)] = time.Now().Add(a.add).Unix()
	a.m.Unlock()
}

func (a *localMap) Set(id any, v []byte) {
	a.set(id, v)
}

func (a *localMap) SetForever(id any, v []byte) {
	a.set(id, v)
}

func (a *localMap) SetJson(id, v any) {
	b, _ := jsoniter.Marshal(v)
	a.set(id, b)
}

func (a *localMap) get(id any) (res any) {
	a.m.RLock()
	res = a.cache[a.sid(id)]
	a.m.RUnlock()
	return
}

func (a *localMap) Get(id any) (res []byte) {
	return cast.Bytes(a.get(id))
}

func (a *localMap) Has(id any) (has bool) {
	return a.Get(id) != nil
}

func (a *localMap) Delete(id any) {
	a.m.Lock()
	delete(a.cache, a.sid(id))
	delete(a.ttl, a.sid(id))
	a.m.Unlock()
}

func (a *localMap) Reset() {
	a.m.Lock()
	a.cache = map[any]any{}
	a.ttl = map[any]int64{}
	a.m.Unlock()
}

func (a *localMap) SetInt(id any, v int) {
	a.set(id, v)
}

func (a *localMap) GetInt(id any) (has bool, v int) {
	p := a.get(id)
	if p == nil {
		return
	}
	v = cast.Int(p)
	return
}

func (a *localMap) SetInts(id any, v []int) {
	a.set(id, v)
}

func (a *localMap) GetInts(id any) (has bool, v []int) {
	p := a.get(id)
	if p == nil {
		return
	}
	v = cast.SliceInt(p)
	return
}

func (a *localMap) deamon() {
	for {
		time.Sleep(time.Second * 20)
		fmt.Println("ttl iterate...")
		now := time.Now().Unix()
		var count int
		a.m.RLock()
		for k, until := range a.ttl {
			if until > now {
				continue
			}
			go a.Delete(k)
			count++
		}
		a.m.RUnlock()
		fmt.Println("ttl deleted", count, "/", len(a.cache))

	}
}

func (a *localMap) Batch(ids ...any) (res []*KV, err error) {

	res = make([]*KV, len(ids))

	for i, x := range ids {
		res[i] = &KV{ID: ids[i], Value: a.Get(x)}
	}

	return

}
