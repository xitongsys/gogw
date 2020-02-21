package common

import (
	"encoding/json"
	"fmt"
)

var UUIDMAP map[string]int = make(map[string]int)

func UUID(key string) string {
	if _, ok := UUIDMAP[key]; !ok {
		UUIDMAP[key] = 0
	}

	UUIDMAP[key]++
	return fmt.Sprint(UUIDMAP[key])
}
