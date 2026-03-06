package cache

import (
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/monopolly/cast"
	"github.com/monopolly/numbers"
	"github.com/niubaoshu/gotiny"
	"github.com/redis/go-redis/v9"
)

type redisEngine struct {
	conn   *redis.Client
	name   string //"users":4242
	prefix string //for iterates
	ttl    time.Duration
}

// localhost:6379
func RedisConn(host string, pass ...string) (conn *redis.Client) {
	p := ""
	if len(pass) > 0 {
		p = pass[0]
	}
	return redis.NewClient(&redis.Options{
		Addr:     host,
		Password: p,
	})
}

func NewRedis(name string, conn *redis.Client, ttl time.Duration) (a Engine) {
	return &redisEngine{
		name:   name,
		prefix: name + ":",
		conn:   conn,
		ttl:    ttl,
	}
}

func NewRedisConn(name string, host, pass string, ttl time.Duration) (a Engine) {
	return &redisEngine{
		conn: RedisConn(host, pass),
		ttl:  ttl,
		name: name,
	}
}

func (a *redisEngine) sid(id any) (res string) {
	return sid(id, a.name) //users:id
}

func (a *redisEngine) set(id any, v []byte) {
	a.conn.Set(context.Background(), a.sid(id), v, a.ttl)
}

func (a *redisEngine) Set(id any, v []byte) {
	a.set(id, v)
}

func (a *redisEngine) SetForever(id any, v []byte) {
	a.conn.Set(context.Background(), a.sid(id), v, 0)
}

func (a *redisEngine) SetJson(id, v any) {
	b, _ := jsoniter.Marshal(v)
	a.set(id, b)
}

func (a *redisEngine) get(id any) (res []byte) {
	p := a.conn.Get(context.Background(), a.sid(id))
	b, er := p.Bytes()
	if er != nil {
		return
	}
	return b
}

func (a *redisEngine) Get(id any) (res []byte) {
	return a.get(id)
}

func (a *redisEngine) Has(id any) (has bool) {
	exist, _ := a.conn.Exists(context.Background(), a.sid(id)).Result()
	return exist == 1
}

func (a *redisEngine) Delete(id any) {
	a.conn.Del(context.Background(), a.sid(id))
}

func (a *redisEngine) SetInt(id any, v int) {
	a.set(id, numbers.IntBytes(v))
}

func (a *redisEngine) GetInt(id any) (has bool, v int) {
	b := a.get(id)
	if b == nil {
		return
	}
	v = numbers.BytesInt(b)
	return
}

func (a *redisEngine) SetInts(id any, v []int) {
	a.set(id, gotiny.Marshal(&v))
}

func (a *redisEngine) GetInts(id any) (has bool, v []int) {
	b := a.get(id)
	if b == nil {
		return
	}
	gotiny.Unmarshal(b, &v)
	return
}

// remove all records users:*
func (a *redisEngine) Reset() {
	var cursor uint64
	var keys []string
	var err error
	ctx := context.Background()
	for {
		// SCAN iterates
		keys, cursor, err = a.conn.Scan(ctx, cursor, a.prefix+"*", 100).Result()
		if err != nil {
			return
		}

		// delete keys
		if len(keys) > 0 {
			a.conn.Del(ctx, keys...).Result()
		}

		// close on 0
		if cursor == 0 {
			break
		}
	}
}

type KV struct {
	ID    any
	Value []byte
}

func (a *redisEngine) Batch(ids ...any) (res []*KV, err error) {

	keys := make([]string, len(ids))
	for i, x := range ids {
		keys[i] = a.sid(x)
	}

	ctx := context.Background()

	vals, err := a.conn.MGet(ctx, keys...).Result()
	if err != nil {
		return
	}

	res = make([]*KV, len(ids))

	for i, v := range vals {
		res[i] = &KV{ID: ids[i], Value: cast.Bytes(v)}
	}

	return res, nil
}
