package cache

//testing
import (
	"testing"
	//testing
	//go test -bench=.
	//go test --timeout 9999999999999s
)

func BenchmarkSID(u *testing.B) {

	u.ReportAllocs()
	for u.Loop() {
		sid(41323, "accounts")
	}
}
