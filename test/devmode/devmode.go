// Package devmode or production
package devmode

import (
	"os"
)

var isDevMode *bool

func IsDevMode() bool {
	return false

	if isDevMode != nil {
		return *isDevMode
	}

	d := os.Getenv("devserver") != ""
	isDevMode = &d

	return *isDevMode
}
