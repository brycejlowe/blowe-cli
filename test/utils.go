package test

import (
	"path/filepath"
	"runtime"
)

func GetRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}
