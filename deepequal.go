package is

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

type equalCheck struct {
	visited map[visit]bool
	err     error
}

type path struct {
	depth int
	path  []string
}

func (p *path) vPath(key reflect.Value) *path {
	depth := p.depth + 1
	pa := append([]string{key.String()}, p.path...)
	return &path{depth: depth, path: pa}

}

func (p *path) sPath(key string) *path {
	depth := p.depth + 1
	pa := append([]string{key}, p.path...)
	return &path{depth: depth, path: pa}

}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (p *path) string() string {
	sb := strings.Builder{}
	first := true
	p.path = reverse(p.path)
	for _, s := range p.path {
		if !first && strings.Index(s, "[") != 0 {
			sb.WriteRune('.')
		}
		sb.WriteString(s)
		first = false
	}
	return sb.String()
}

func (p *path) arrayPath(index int) *path {
	depth := p.depth + 1
	pa := append([]string{fmt.Sprintf("[%d]", index)}, p.path...)
	return &path{depth: depth, path: pa}
}

func (eq *equalCheck) unsafeEqual(v1, v2 reflect.Value, path *path) (isEq bool) {
	var uv1, uv2 interface{}
	switch v1.Kind() {
	case reflect.Bool:
		uv1 = v1.Bool()
		uv2 = v2.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		uv1 = v1.Int()
		uv2 = v2.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		uv1 = v1.Uint()
		uv2 = v2.Uint()
	case reflect.Float32, reflect.Float64:
		uv1 = v1.Float()
		uv2 = v2.Float()
	case reflect.Complex64, reflect.Complex128:
		uv1 = v1.Complex()
		uv2 = v2.Complex()
	case reflect.String:
		uv1 = v1.String()
		uv2 = v2.String()
	case reflect.UnsafePointer:
		uv1 = v1.Pointer()
		uv2 = v2.Pointer()
	}

	if uv1 != uv2 {
		eq.err = fmt.Errorf("Values are not equal\n%s: '%v' is not equal to '%v'", path.string(), uv1, uv2)
		return false
	}
	return true
}

func isNil(v reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func (eq *equalCheck) visit(v1, v2 reflect.Value) bool {
	switch v1.Kind() {
	case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
		if v1.CanAddr() && v2.CanAddr() {
			addr1 := unsafe.Pointer(v1.UnsafeAddr())
			addr2 := unsafe.Pointer(v2.UnsafeAddr())
			if uintptr(addr1) > uintptr(addr2) {
				addr1, addr2 = addr2, addr1
			}
			typ := v1.Type()
			v := visit{addr1, addr2, typ}
			if eq.visited[v] {
				return true
			}

			eq.visited[v] = true
			return false
		}
	}
	return false

}
func (eq *equalCheck) deepValueEqual(v1, v2 reflect.Value, path *path) bool {
	if !v1.IsValid() || !v2.IsValid() {
		if v1.IsValid() != v2.IsValid() {
			eq.err = fmt.Errorf("Values are not equal\n%s: One value is zero (untyped nil)", path.string())
			return false
		}
		return true
	}
	if v1.Type() != v2.Type() {
		eq.err = fmt.Errorf("Values are not equal\n%s: Different types '%s' and '%s'", path.string(), v1.Type(), v2.Type())
		return false
	}
	if isNil(v1) != isNil(v2) {
		eq.err = fmt.Errorf("Values are not equal\n%s: One value is nil", path.string())
		return false
	}

	if isNil(v1) == true {
		return true
	}

	if eq.visit(v1, v2) {
		return true
	}

	switch v1.Kind() {
	case reflect.Array:
		if v1.Len() != v2.Len() {
			eq.err = fmt.Errorf("Values are not equal\n%s: Array length is different", path.string())
			return false
		}
		for i := 0; i < v1.Len(); i++ {
			if !eq.deepValueEqual(v1.Index(i), v2.Index(i), path.arrayPath(i)) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if v1.Len() != v2.Len() {
			eq.err = fmt.Errorf("Values are not equal\n%s: Slice length is different", path.string())
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			if !eq.deepValueEqual(v1.Index(i), v2.Index(i), path.arrayPath(i)) {
				return false
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() {
			return true
		}
		return eq.deepValueEqual(v1.Elem(), v2.Elem(), path)
	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		return eq.deepValueEqual(v1.Elem(), v2.Elem(), path)
	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if v := v1.Type().Field(i).Tag.Get("is"); v == "-" {
				continue
			}
			if !eq.deepValueEqual(v1.Field(i), v2.Field(i), path.sPath(v1.Type().Field(i).Name)) {
				return false
			}
		}
		return true
	case reflect.Map:
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() || !eq.deepValueEqual(val1, val2, path.vPath(k)) {
				return false
			}
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		return false
	default:
		return eq.unsafeEqual(v1, v2, path)
	}
}

func deepEqual(x, y interface{}) (isEq bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			isEq = false
			err = fmt.Errorf("panic: %s", r)
		}
	}()

	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)

	if !v1.IsValid() || !v2.IsValid() {
		if v1.IsValid() != v2.IsValid() {
			err = fmt.Errorf("Values are not equal\nOne value is zero (untyped nil)")
			return false, err
		}
		return true, nil
	}
	if v1.Type() != v2.Type() {
		err = fmt.Errorf("Values are not equal\nDifferent types '%s' and '%s'", v1.Type(), v2.Type())
		return false, err
	}

	eq := equalCheck{make(map[visit]bool), nil}

	isEq = eq.deepValueEqual(v1, v2, &path{depth: 0, path: name(v1)})
	err = eq.err
	return
}

func name(v1 reflect.Value) []string {
	getName := func() string {
		name := v1.Type().Name()
		if name == "" {
			name = v1.Kind().String()
		}
		return name
	}

	switch v1.Kind() {
	case reflect.Struct, reflect.Interface:
		return []string{getName()}
	case reflect.Ptr, reflect.UnsafePointer:
		return name(v1.Elem())

	}
	return []string{}
}
