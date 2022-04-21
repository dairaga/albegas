package app

import (
	"fmt"
	"reflect"
	"strings"
)

type models map[string]reflect.Value

func (m models) String() string {
	var info []string

	for k, v := range m {
		info = append(info, fmt.Sprintf("%s:%s", k, v.Type().String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(info, ", "))
}

// -----------------------------------------------------------------------------

func (m models) Get(name string) (interface{}, bool) {
	return retcheck(m[name])
}

// -----------------------------------------------------------------------------

func (m models) Index(name string, idx int) (interface{}, bool) {
	val, ok := m[name]

	if ok && valIsSlice(val) {
		return retcheck(val.Index(idx))
	}
	return nil, false
}

// -----------------------------------------------------------------------------

func (m models) MapIndex(name string, key interface{}) (interface{}, bool) {
	val, ok := m[name]

	if ok && valIsMap(val) {
		return retcheck(val.MapIndex(reflect.ValueOf(key)))
	}

	return nil, false
}

// -----------------------------------------------------------------------------

func (m models) Set(name string, val interface{}) {
	x, ok := m[name]
	if ok {
		x.Set(reflect.ValueOf(val))
	} else {
		x = reflect.ValueOf(val)
	}

	m[name] = x
}

// -----------------------------------------------------------------------------

func (m models) Append(name string, val1 interface{}, vals ...interface{}) {
	args := make([]reflect.Value, 1+len(vals))
	args[0] = reflect.ValueOf(val1)

	for i := range vals {
		args[i+1] = reflect.ValueOf(vals[i])
	}

	x, ok := m[name]
	if !ok || !valIsSlice(x) {
		x = reflect.MakeSlice(reflect.SliceOf(args[0].Type()), 0, 0)
	}

	x = reflect.Append(x, args...)

	m[name] = x
}

// -----------------------------------------------------------------------------

func (m models) SetIndex(name string, idx int, val interface{}) {
	x, ok := m[name]
	if ok && valIsSlice(x) {
		x.Index(idx).Set(reflect.ValueOf(val))
		m[name] = x
	} else {
		m.Append(name, val)
	}
}

// -----------------------------------------------------------------------------

func (m models) SetMapIndex(name string, key, val interface{}) {
	x, ok := m[name]
	kval := reflect.ValueOf(key)
	vval := reflect.ValueOf(val)

	if !ok || !valIsMap(x) {
		x = reflect.MakeMap(reflect.MapOf(kval.Type(), vval.Type()))
	}

	x.MapIndex(kval).Set(vval)
	m[name] = x
}

// -----------------------------------------------------------------------------

func retcheck(val reflect.Value) (interface{}, bool) {
	if val.IsValid() {
		return val.Interface(), true
	}

	return nil, false
}

// -----------------------------------------------------------------------------

func valIsSlice(val reflect.Value) bool {
	if val.IsValid() {
		kind := val.Kind()
		return reflect.Slice == kind || reflect.Array == kind
	}

	return false
}

// -----------------------------------------------------------------------------

func valIsMap(val reflect.Value) bool {
	return val.IsValid() && val.Kind() == reflect.Map
}
