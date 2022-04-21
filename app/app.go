//go:build js && wasm
// +build js,wasm

package app

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"unsafe"

	"time"
)

const Self = "self"
const sep = "::"

type app struct {
	mutex      *sync.Mutex
	models     map[string]*model
	components map[string]*component
	m2v        map[string]map[string]bool // model 2 component
	v2m        map[string]map[string]bool // component 2 model
	mvvm       map[string]reflect.Value
	ch         chan struct{}
}

var _app = &app{
	mutex:      new(sync.Mutex),
	models:     make(map[string]*model),
	components: make(map[string]*component),

	ch: make(chan struct{}),
}

// -----------------------------------------------------------------------------

func (a *app) append(tatto string, com *component) {
	a.mutex.Lock()
	a.components[tatto] = com
	a.mutex.Unlock()
}

// -----------------------------------------------------------------------------

func (a *app) remove(tatto string) {
	a.mutex.Lock()
	ms, ok := a.v2m[tatto]
	if ok && ms != nil {
		for model, v := range ms {
			if v {
				delete(a.m2v[model], tatto)
				delete(a.mvvm, typKey(model, tatto))
			}
		}
		delete(a.v2m, tatto)
	}
	delete(a.components, tatto)
	a.mutex.Unlock()
}

// -----------------------------------------------------------------------------

func (a *app) set(name string, value interface{}) {
	x, ok := a.models[name]
	if ok && x != nil {
		x.value = value
		return
	}

	a.models[name] = &model{
		name:  name,
		value: value,
	}
}

// -----------------------------------------------------------------------------

func (a *app) invoke(model, componet string, value interface{}) {
	fn, ok := a.mvvm[typKey(model, componet)]
	if ok && fn.IsValid() {
		com := a.components[componet]
		fn.Call(
			[]reflect.Value{
				reflect.ValueOf(com),
				reflect.ValueOf(value),
			},
		)
	}
}

// -----------------------------------------------------------------------------

func (a *app) trigger(model string, value interface{}) {
	components := a.m2v[model]
	if components != nil {
		for c, ok := range components {
			if ok {
				a.invoke(model, c, value)
			}
		}
	}
}

// -----------------------------------------------------------------------------

func Set(name string, value interface{}) {
	_app.set(name, value)
	_app.trigger(name, value)
}

// -----------------------------------------------------------------------------

func Append(name string, value interface{}) {
	x, ok := _app.models[name]
	if ok && x != nil {
		x.value = reflect.Append(reflect.ValueOf(x.value), reflect.ValueOf(value)).Interface()
		_app.trigger(name, x.value)
		return
	}
	fmt.Printf("warn: data [%s] not found\n", name)
}

// -----------------------------------------------------------------------------

func SetIndex(name string, idx int, value interface{}) {
	x, ok := _app.models[name]

	if ok && x != nil {
		reflect.ValueOf(x.value).Index(idx).Set(reflect.ValueOf(value))
		_app.trigger(name, x.value)
		_app.trigger(fmt.Sprintf("%s[%d]", name, idx), value)
	}
	fmt.Printf("warn: data [%s] not found\n", name)
}

// -----------------------------------------------------------------------------

func SetMapIndex(name string, key interface{}, value interface{}) {
	x, ok := _app.models[name]

	if ok && x != nil {
		reflect.ValueOf(x.value).MapIndex(reflect.ValueOf(key)).Set(reflect.ValueOf(value))
		_app.trigger(name, x.value)
		_app.trigger(fmt.Sprintf("%s[%v]", name, key), value)
	}
	fmt.Printf("warn: data [%s] not found\n", name)
}

// -----------------------------------------------------------------------------

func Get(name string) (interface{}, bool) {
	x, ok := _app.models[name]
	if ok && x != nil {
		return x.value, true
	}
	return nil, false
}

// -----------------------------------------------------------------------------

func Index(name string, idx int) (interface{}, bool) {
	x, ok := _app.models[name]
	if ok && x != nil {
		return any(reflect.ValueOf(x.value).Index(idx))
	}

	return nil, false
}

// -----------------------------------------------------------------------------

func MapIndex(name string, key interface{}) (interface{}, bool) {
	x, ok := _app.models[name]
	if ok && x != nil {
		return any(
			reflect.ValueOf(x.value).MapIndex(reflect.ValueOf(key)),
		)

	}

	return nil, false
}

// -----------------------------------------------------------------------------

func Run() {
	<-_app.ch

	for _, com := range _app.components {
		com.Depose()
	}
}

// -----------------------------------------------------------------------------

func Terminal() {
	_app.ch <- struct{}{}
}

// -----------------------------------------------------------------------------

func typKey(a, b string) string {
	return a + sep + b
}

// -----------------------------------------------------------------------------

func any(val reflect.Value) (interface{}, bool) {
	if val.IsValid() {
		return val.Interface(), true
	}
	return nil, false
}

// -----------------------------------------------------------------------------

/*
	Reference from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
*/

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	randSrc = rand.NewSource(time.Now().UnixNano())
)

func tatto(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
