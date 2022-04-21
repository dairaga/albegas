//go:build js && wasm
// +build js,wasm

package app

import (
	"fmt"
	"reflect"
	"sync"
)

const Self = "self"
const sep = "::"

type app struct {
	mutex      *sync.Mutex
	models     models
	components map[string]*component
	mvvm       *mvvm
	ch         chan struct{}
}

var _app = &app{
	mutex:      new(sync.Mutex),
	models:     make(map[string]reflect.Value),
	components: make(map[string]*component),
	mvvm:       newMVVM(),

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
	a.mvvm.Unbind(tatto)
	delete(a.components, tatto)
	a.mutex.Unlock()
}

// -----------------------------------------------------------------------------

func (a *app) trigger(model string, value interface{}) {
	tattos := a.mvvm.Views(model)
	for tatto := range tattos {
		com := a.components[tatto]
		if com != nil {
			a.mvvm.trigger(model, tatto, com, value)
		}
	}
}

// -----------------------------------------------------------------------------

func (a *app) triggerIndx(model string, idx int, value interface{}) {
	a.trigger(fmt.Sprintf("%s[%d]", model, idx), value)
}

// -----------------------------------------------------------------------------

func (a *app) triggerMapIndex(model string, key, val interface{}) {
	a.trigger(fmt.Sprintf("%s[%v]", model, key), val)
}

// -----------------------------------------------------------------------------

func Set(name string, value interface{}) {
	_app.models.Set(name, value)
	_app.trigger(name, value)
}

// -----------------------------------------------------------------------------

func Append(name string, value interface{}) {
	_app.models.Append(name, value)
	_app.trigger(name, value)
}

// -----------------------------------------------------------------------------

func SetIndex(name string, idx int, value interface{}) {
	_app.models.SetIndex(name, idx, value)
	_app.triggerIndx(name, idx, value)
}

// -----------------------------------------------------------------------------

func SetMapIndex(name string, key interface{}, value interface{}) {
	_app.models.SetMapIndex(name, key, value)
	_app.triggerMapIndex(name, key, value)
}

// -----------------------------------------------------------------------------

func Get(name string) (interface{}, bool) {
	return _app.models.Get(name)
}

// -----------------------------------------------------------------------------

func Index(name string, idx int) (interface{}, bool) {
	return _app.models.Index(name, idx)
}

// -----------------------------------------------------------------------------

func MapIndex(name string, key interface{}) (interface{}, bool) {
	return _app.models.MapIndex(name, key)
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
