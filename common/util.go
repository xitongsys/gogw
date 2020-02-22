package common

import (
	"fmt"
	"sync"
)

var mu = sync.Mutex{}

var UUIDMAP map[string]int = make(map[string]int)

func UUID(key string) string {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := UUIDMAP[key]; !ok {
		UUIDMAP[key] = 0
	}

	UUIDMAP[key]++
	return fmt.Sprint(UUIDMAP[key])
}

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func Min(a, b int) int {
	if a > b {
		return b
	}

	return a
}