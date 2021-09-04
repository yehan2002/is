//go:build !go1.17
// +build !go1.17

package is

import "reflect"

func isExported(method reflect.Method) bool { return method.PkgPath == "" }
