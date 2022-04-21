package app

import (
	"fmt"
	"reflect"
	"strings"
)

type mvvm struct {
	m2v map[string]map[string]bool
	v2m map[string]map[string]bool

	funcs map[string]reflect.Value
}

// -----------------------------------------------------------------------------

func (x *mvvm) String() string {
	var info []string
	for k := range x.funcs {
		info = append(info, k)
	}
	return fmt.Sprintf("[%s]", strings.Join(info, ", "))
}

// -----------------------------------------------------------------------------

func (x *mvvm) rel(a, b string) string {
	return a + "::" + b
}

// -----------------------------------------------------------------------------

func (x *mvvm) add(table map[string]map[string]bool, key, val string) {
	sub := table[key]
	if sub == nil {
		sub = make(map[string]bool)
		table[key] = sub
	}
	sub[val] = true
}

// -----------------------------------------------------------------------------

func (x mvvm) Views(model string) map[string]bool {
	return x.m2v[model]
}

// -----------------------------------------------------------------------------

func (x mvvm) Models(view string) map[string]bool {
	return x.v2m[view]
}

// -----------------------------------------------------------------------------

func (x *mvvm) Watch(m, v string, fn interface{}) {
	// add relation to m2v
	x.add(x.m2v, m, v)

	// add releation to v2m
	x.add(x.v2m, v, m)

	x.funcs[x.rel(m, v)] = reflect.ValueOf(fn)
}

// -----------------------------------------------------------------------------

func (x *mvvm) WatchIndex(m string, idx int, v string, fn interface{}) {
	x.Watch(fmt.Sprintf("%s[%d]", m, idx), v, fn)
}

// -----------------------------------------------------------------------------

func (x *mvvm) WatchMapIndex(m string, key interface{}, v string, fn interface{}) {
	x.Watch(fmt.Sprintf("%s[%v]", m, key), v, fn)
}

// -----------------------------------------------------------------------------

func (x *mvvm) Unwatch(m, v string) {
	delete(x.m2v[m], v)
	delete(x.v2m[v], m)
	delete(x.funcs, x.rel(m, v))
}

// -----------------------------------------------------------------------------

func (x *mvvm) UnwatchIndex(m string, idx int, v string) {
	x.Unwatch(fmt.Sprintf("%s[%d]", m, idx), v)
}

// -----------------------------------------------------------------------------

func (x *mvvm) UnwatchMapIndex(m string, key interface{}, v string) {
	x.Unwatch(fmt.Sprintf("%s[%v]", m, key), v)
}

// -----------------------------------------------------------------------------

func (x *mvvm) Unbind(v string) {
	ms := x.v2m[v]

	for m := range ms {
		x.Unwatch(m, v)
	}

	delete(x.v2m, v)
}

// -----------------------------------------------------------------------------

func (x *mvvm) trigger(m, v string, a, b interface{}) {
	fn, ok := x.funcs[x.rel(m, v)]
	if ok && fn.IsValid() && fn.Type().Kind() == reflect.Func {
		fn.Call([]reflect.Value{
			reflect.ValueOf(a),
			reflect.ValueOf(b),
		})
	}
}

// -----------------------------------------------------------------------------

func (x *mvvm) triggerIndex(m string, idx int, v string, a, b interface{}) {
	x.trigger(fmt.Sprintf("%s[%d]", m, idx), v, a, b)
}

// -----------------------------------------------------------------------------

func (x *mvvm) triggerMapIndex(m string, key interface{}, v string, a, b interface{}) {
	x.trigger(fmt.Sprintf("%s[%v]", m, key), v, a, b)
}

// -----------------------------------------------------------------------------

func newMVVM() *mvvm {
	return &mvvm{
		m2v:   make(map[string]map[string]bool),
		v2m:   make(map[string]map[string]bool),
		funcs: make(map[string]reflect.Value),
	}
}
