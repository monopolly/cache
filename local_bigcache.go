package cache

import (
	"context"
	"encoding/binary"
	"log"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/niubaoshu/gotiny"
)

type localCache struct {
	cache  *bigcache.BigCache
	name   string
	prefix string
}

// big cache
func NewLocal(name string, clean, ttl time.Duration) (res Engine) {
	defer log.Printf("local cache: %s\n", name)
	var c bigcache.Config
	c.CleanWindow = clean
	c.LifeWindow = ttl
	c.Shards = 512
	c.Verbose = false
	cache, _ := bigcache.New(context.Background(), c)

	p := new(localCache)
	p.cache = cache
	p.name = name
	p.prefix = name + ":"
	return p
}

func (a *localCache) Set(id any, v []byte) {
	a.cache.Set(a.sid(id), v)
}

func (a *localCache) SetForever(id any, v []byte) {
	a.Set(id, v)
}

func (a *localCache) SetJson(id, v any) {
	b, _ := jsoniter.Marshal(v)
	a.cache.Set(a.sid(id), b)
}

func (a *localCache) Get(id any) (res []byte) {
	res, _ = a.cache.Get(a.sid(id))
	return
}

func (a *localCache) Has(id any) (has bool) {
	return a.Get(id) != nil
}

func (a *localCache) Delete(id any) {
	a.cache.Delete(a.sid(id))
}

func (a *localCache) Reset() {
	p := a.cache.Iterator()
	for p.SetNext() {
		v, er := p.Value()
		if er != nil {
			continue
		}
		if strings.HasPrefix(v.Key(), a.prefix) {
			a.Delete(v.Key())
		}
	}
}

func (a *localCache) SetInt(id any, v int) {
	a.cache.Set(a.sid(id), IntBytes(v))
}

func (a *localCache) GetInt(id any) (v int) {
	b, _ := a.cache.Get(a.sid(id))
	if b == nil {
		return
	}
	return BytesInt(b)
}

func (a *localCache) SetInts(id any, v []int) {
	a.cache.Set(a.sid(id), gotiny.Marshal(&v))
}

func (a *localCache) GetInts(id any) (v []int) {
	b, _ := a.cache.Get(a.sid(id))
	if b == nil {
		return
	}
	gotiny.Unmarshal(b, &v)
	return
}

func (a *localCache) sid(id any) (res string) {
	return sid(id, a.name)
}

// int
func IntBytes(i int) (r []byte) {
	r = make([]byte, 8)
	binary.LittleEndian.PutUint64(r, uint64(i))
	return
}

func BytesInt(b []byte) (i int) {
	if b == nil {
		return 0
	}
	return int(binary.LittleEndian.Uint64(b))
}
