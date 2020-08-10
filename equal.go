package is

import (
	"fmt"
	"math"
	"reflect"
)

const (
	eqNil = iota
	eqValZero
	eqChannel
	eqNotEqual
	eqDiffTypes
	eqDiffLenArray
	eqDiffLenSlice
	eqIncomparable
)

const largeList = 20

type visit struct {
	a1  uintptr
	a2  uintptr
	typ reflect.Type
}

func compare(name string, v1, v2 reflect.Value, visited map[visit]bool) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			if eq, ok := r.(*eqError); ok {
				if name != "" {
					if eq.path != "" && eq.path[0] != '[' {
						eq.path += "."
					}
					eq.path = name + eq.path
				}
			}
			panic(r)
		}
	}()

	if !v1.IsValid() || !v2.IsValid() {
		if v1.IsValid() || v2.IsValid() {
			panic(equalityError(eqValZero))
		}
		return true // both values are untyped nil
	}

	if v1.Type() != v2.Type() {
		panic(equalityError(eqDiffTypes, v1.Type(), v2.Type()))
	}

	switch v1.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.String, reflect.Uintptr:
		if v1.CanInterface() {
			if v1.Interface() != v2.Interface() {
				panic(equalityError(eqNotEqual, v1.Interface(), v2.Interface()))
			}
		} else {
			switch v1.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				if v1.Uint() != v2.Uint() {
					panic(equalityError(eqNotEqual, v1.Uint(), v2.Uint()))
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v1.Uint() != v2.Uint() {
					panic(equalityError(eqNotEqual, v1.Int(), v2.Int()))
				}
			case reflect.Bool:
				if v1.Bool() != v2.Bool() {
					panic(equalityError(eqNotEqual, v1.Bool(), v2.Bool()))
				}
			case reflect.String:
				if v1.String() != v2.String() {
					panic(equalityError(eqNotEqual, v1.String(), v2.String()))
				}
			}
		}
		return true
	case reflect.Float32, reflect.Float64:
		if math.Float64bits(v1.Float()) != math.Float64bits(v2.Float()) {
			panic(equalityError(eqNotEqual, v1.Float(), v2.Float()))
		}
		return true
	case reflect.Complex64, reflect.Complex128:
		c1, c2 := v1.Complex(), v2.Complex()
		if math.Float64bits(real(c1)) != math.Float64bits(real(c2)) || math.Float64bits(imag(c1)) != math.Float64bits(imag(c2)) {
			panic(equalityError(eqNotEqual, c1, c2))
		}
		return true
	case reflect.Chan:
		if v1.Pointer() != v2.Pointer() {
			panic(equalityError(eqChannel))
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		panic(equalityError(eqIncomparable, reflect.Func))
	case reflect.UnsafePointer:
		if v1.IsNil() != v2.IsNil() {
			panic(equalityError(eqNil))
		}
		if v1.IsNil() {
			return true //both are nil
		}
		if v1.Pointer() != v2.Pointer() {
			panic(equalityError(eqNotEqual, v1.Pointer(), v2.Pointer()))
		}
		return true

	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if v := v1.Type().Field(i).Tag.Get("is"); v == "-" {
				continue
			}
			compare(v1.Type().Field(i).Name, v1.Field(i), v2.Field(i), visited)
		}
		return true
	case reflect.Array:
		if v1.Len() != v2.Len() {
			panic(equalityError(eqDiffLenArray, v1.Len(), v2.Len()))
		}
		if v1.CanInterface() && v1.Type().Comparable() && v1.Interface() == v2.Interface() {
			return true // array is equal
		}
		if v1.Len() > largeList {
			for i := 0; i < v1.Len(); i++ {
				compare("", v1.Index(i), v2.Index(i), visited)
			}
		} else {
			for i := 0; i < v1.Len(); i++ {
				compare(fmt.Sprintf("[%d]", i), v1.Index(i), v2.Index(i), visited)
			}
		}
		return true
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			panic(equalityError(eqNil))
		}
		if v1.IsNil() {
			return true //both are nil
		}
		if v1.Pointer() == v2.Pointer() {
			return true //same memory address
		}

		if _, ok := visited[visit{v1.Pointer(), v2.Pointer(), v1.Type()}]; ok {
			return true //already compared
		}
		visited[visit{v1.Pointer(), v2.Pointer(), v1.Type()}] = true

		if v1.Len() != v2.Len() {
			panic(equalityError(eqDiffLenSlice, v1.Len(), v2.Len()))
		}
		if v1.CanInterface() && v1.Type().Comparable() && v1.Interface() == v2.Interface() {
			return true // array is equal
		}
		if v1.Len() > largeList {
			for i := 0; i < v1.Len(); i++ {
				compare("", v1.Index(i), v2.Index(i), visited)
			}
		} else {
			for i := 0; i < v1.Len(); i++ {
				compare(fmt.Sprintf("[%d]", i), v1.Index(i), v2.Index(i), visited)
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() != v2.IsNil() {
			panic(equalityError(eqNil))
		}
		if v1.CanAddr() {
			if _, ok := visited[visit{v1.UnsafeAddr(), v2.UnsafeAddr(), v1.Type()}]; ok {
				return true //already compared
			}
			visited[visit{v1.UnsafeAddr(), v2.UnsafeAddr(), v1.Type()}] = true
		}
		if v1.IsNil() {
			return true //both are nil
		}
		return compare("", v1.Elem(), v2.Elem(), visited)
	case reflect.Ptr:
		if v1.IsNil() != v2.IsNil() {
			panic(equalityError(eqNil))
		}
		if v1.IsNil() {
			return true //both are nil
		}
		if v1.Pointer() == v2.Pointer() {
			return true //same memory address
		}
		if _, ok := visited[visit{v1.Pointer(), v2.Pointer(), v1.Type()}]; ok {
			return true //already compared
		}
		visited[visit{v1.Pointer(), v2.Pointer(), v1.Type()}] = true
		return compare("", v1.Elem(), v2.Elem(), visited)
	case reflect.Map:
		if v1.Len() != v2.Len() {
			panic(equalityError(eqNil))
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		if _, ok := visited[visit{v1.Pointer(), v2.Pointer(), v1.Type()}]; ok {
			return true //already compared
		}
		visited[visit{v1.Pointer(), v2.Pointer(), v1.Type()}] = true
		for _, k := range v1.MapKeys() {
			compare(fmt.Sprint(k), v1.MapIndex(k), v2.MapIndex(k), visited)
		}
		return true
	default:
		panic("unreachable")
	}
}

func isEqual(v1, v2 interface{}) (isEq bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if eq, ok := r.(*eqError); ok {
				isEq = false
				err = eq
				return
			}
			panic(r)
		}
		isEq = true
	}()
	compare("", reflect.ValueOf(v1), reflect.ValueOf(v2), make(map[visit]bool))
	return
}

type eqError struct {
	err  string
	path string
}

func equalityError(err int, v ...interface{}) error {
	var s string
	switch err {
	case eqDiffTypes:
		s = "Different types '%s' and '%s'"
	case eqValZero:
		s = "One value is invalid (untyped nil)"
	case eqIncomparable:
		s = "Incomparable kind: %s"
	case eqNotEqual:
		s = "'%v' is not equal to '%v'"
	case eqDiffLenArray:
		s = "Array length is different (%d,%d)"
	case eqDiffLenSlice:
		s = "Slice length is different (%d,%d)"
	case eqChannel:
		s = "Channels reffer to different memory locations."
	case eqNil:
		s = "One value is nil"
	default:
		panic("unknown error")
	}
	return &eqError{err: fmt.Sprintf(s, v...)}
}

func (e *eqError) Error() string {
	return fmt.Sprintf("Values are not equal:\n%s: %s", e.path, e.err)
}
