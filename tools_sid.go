package cache

import (
	"fmt"
	"strconv"
)

func sid(id any, name string) (res string) {
	switch k := id.(type) {
	case string:
		res = k
	case int:
		res = strconv.Itoa(k)
	case []byte:
		res = string(k)
	default:
		res = fmt.Sprint(id)
	}

	return fmt.Sprintf("%s:%s", name, res)
}
