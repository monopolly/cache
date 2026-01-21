package cache

import (
	"fmt"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type localMap struct {
	cache map[any]any
	ttl   map[any]int64
	m     sync.RWMutex
	max   int
	add   time.Duration
}

// simple map
func NewMap(ttl time.Duration) (res Engine) {
	a := new(localMap)
	a.cache = map[any]any{}
	a.ttl = map[any]int64{}
	a.add = ttl
	go a.deamon()
	return a
}

func (a *localMap) set(id any, v any) {
	a.m.Lock()
	a.cache[id] = v
	a.ttl[id] = time.Now().Add(a.add).Unix()
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
	res = a.cache[id]
	a.m.RUnlock()
	return
}

func (a *localMap) Get(id any) (res []byte) {
	res, _ = a.get(id).([]byte)
	return
}

func (a *localMap) Has(id any) (has bool) {
	return a.Get(id) != nil
}

func (a *localMap) Delete(id any) {
	a.m.Lock()
	delete(a.cache, id)
	delete(a.ttl, id)
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

func (a *localMap) GetInt(id any) (v int) {
	v, _ = a.get(id).(int)
	return
}

func (a *localMap) SetInts(id any, v []int) {
	a.set(id, v)
}

func (a *localMap) GetInts(id any) (v []int) {
	v, _ = a.get(id).([]int)
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
