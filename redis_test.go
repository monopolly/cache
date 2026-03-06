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
	_, v := p.GetInt("t1")
	if v != 42 {
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

	p.Set(1, []byte("1"))
	p.Set(2, []byte("2"))
	p.Set(3, []byte("3"))
	p.Set(4, []byte("4"))
	p.Set(5, []byte("5"))

	list, err := p.Batch(1, 2, 3, 4, 5)
	if err != nil {
		panic(err)
	}

	for _, x := range list {
		fmt.Println(x.ID, string(x.Value))
	}

	// select {}
}

func TestRedisLocalDB(u *testing.T) {
	___(u)

	p := NewRedisConn("333", localhost, localpass, time.Minute)

	p.SetInt("t1", 42)
	_, v := p.GetInt("t1")
	if v != 42 {
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
