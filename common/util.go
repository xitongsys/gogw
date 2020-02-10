package common

import (
	"github.com/satori/go.uuid"
)

func uuid() (string, error) {
	u2, err := uuid.NewV4()
	return u2.String(), err
}
