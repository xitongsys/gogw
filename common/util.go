package common

import (
	"fmt"
	
	"github.com/google/uuid"
)

var ID = 0

func UUID() string {
	ID++
	return fmt.Sprint(ID)
}

func UUID0() string {
	return uuid.New().String()
}
