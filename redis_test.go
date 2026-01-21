package cache

//testing
import (
	"fmt"
	"strings"
	"testing"
	"time"
	//testing
	//go test -bench=.
	//go test --timeout 9999999999999s
)

const (
	localhost = "localhost:6379"
	localpass = ""
)

func TestRedisLocal(u *testing.T) {
	___(u)

	conn := RedisConn(localhost, localpass)

	p := NewRedis("test", conn, time.Minute)

	p.SetInt("t1", 42)
	if p.GetInt("t1") != 42 {
		panic("error")
	}

	has := p.Has(11)
	if has {
		panic("has error")
	}

	if !p.Has("t1") {
		panic("non has error")
	}

	fmt.Println(sid(1, "user"))
}

func TestRedisLocalDB(u *testing.T) {
	___(u)

	p := NewRedisConn("333", localhost, localpass, time.Minute)

	p.SetInt("t1", 42)
	if p.GetInt("t1") != 42 {
		panic("none equal")
	}

	has := p.Has(11)
	if has {
		panic("has error")
	}

	if !p.Has("t1") {
		panic("non has error")
	}
}

func ___(u *testing.T) {
	fmt.Printf("\033[1;32m%s\033[0m\n", strings.ReplaceAll(u.Name(), "Test", ""))
}

func BenchmarkRedisLocal(u *testing.B) {

	conn := RedisConn(localhost, localpass)
	p := NewRedis("test", conn, time.Minute)
	u.ReportAllocs()
	for u.Loop() {
		p.Get("t1")
	}
}

func Benchmark2Local(u *testing.B) {

	conn := RedisConn(localhost, localpass)
	p := NewRedis("test", conn, time.Minute)
	u.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Get("t1")
		}
	})
}
