package cache

//testing
import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	//testing
	//go test -bench=.
	//go test --timeout 9999999999999s
)

func TestMap(u *testing.T) {
	___(u)

	p := NewMap(time.Minute)
	p.SetInt(1, 42)
	if p.GetInt(1) != 42 {
		panic("non 42")
	}
	p.SetInts(1, []int{42, 43})
	fmt.Println(p.GetInts(1))

	p.Set(1, []byte("42"))
	fmt.Println(string(p.Get(1)))

	b, _ := json.Marshal(map[string]any{
		"name":  "James Miller",
		"names": "James Miller",
		"nam":   "James Miller",
		"na":    "James Miller",
		"n":     "James Miller",
		"a1":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		"a2":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		"a3":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		"a4":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		"a5":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		"a6":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		"a7":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		"a8":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		"b1":    true,
		"b2":    true,
		"b3":    true,
		"b4":    true,
		"b5":    true,
		"b6":    true,
	})

	fmt.Println("len load", len(b))
	go func() {
		for x := range 100000000000 {
			p.Set(x, b)
			time.Sleep(time.Millisecond)
		}
	}()

	go func() {
		for x := range 100000000000 {
			p.Set(x, b)
			time.Sleep(time.Millisecond)
		}
	}()

	select {}

}

func BenchmarkMap(u *testing.B) {
	p := NewMap(time.Hour)
	p.SetInt(1, 42)
	u.ReportAllocs()
	for u.Loop() {
		p.GetInt(1)
	}
}

func BenchmarkMap2(u *testing.B) {
	p := NewMap(time.Hour)
	p.SetInt(1, 42)
	u.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.GetInt(1)
		}
	})
}
