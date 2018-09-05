package utils

import (
	"reflect"
	"runtime"
	"strings"
)

// FuncName gets the name of the function at runtime
func FuncName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	parts := strings.Split(name, ".")
	name = parts[len(parts)-2]
	return name
}
