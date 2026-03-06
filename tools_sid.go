package cache

import (
	"fmt"
	"strconv"
)

func sid(id any, name string) string {
	switch v := id.(type) {
	case string:
		return name + ":" + v
	case int:
		return name + ":" + strconv.Itoa(v)
	case int64:
		return name + ":" + strconv.FormatInt(v, 10)
	case uint64:
		return name + ":" + strconv.FormatUint(v, 10)
	case []byte:
		return name + ":" + string(v)
	default:
		return name + ":" + fmt.Sprint(v)
	}
}
