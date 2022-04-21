//go:build js && wasm
// +build js,wasm

package app

import (
	"fmt"
	"syscall/js"

	"github.com/dairaga/albegas"
	"github.com/dairaga/albegas/element"
	"github.com/dairaga/albegas/global"
)

type component struct {
	parent *component
	tatto  string
	albegas.Element
	callbacks map[string]js.Func
}

// -----------------------------------------------------------------------------

func (c *component) target(elm string) albegas.Element {
	if Self == elm {
		return c.Element
	}
	return c.Query(elm)
}

// -----------------------------------------------------------------------------

func (c *component) on(elm, typ string, listener js.Func) {
	target := c.target(elm)
	target.AddEventListener(typ, listener)
	c.callbacks[typKey(elm, typ)] = listener
}

// -----------------------------------------------------------------------------

func (c *component) Append(elm string, cmp albegas.Component) {
	c.target(elm).Ref().Call("append", cmp.Ref())
}

// -----------------------------------------------------------------------------

func (c *component) Watch(model string, fn interface{}) {
	_app.mvvm.Watch(model, c.tatto, fn)
}

// -----------------------------------------------------------------------------

func (c *component) WatchIndex(name string, idx int, fn interface{}) {
	_app.mvvm.WatchIndex(name, idx, c.tatto, fn)
}

// -----------------------------------------------------------------------------

func (c *component) WatchMapIndex(name string, key interface{}, fn interface{}) {
	_app.mvvm.WatchMapIndex(name, key, c.tatto, fn)
}

// -----------------------------------------------------------------------------

func (c *component) On(elm, typ string, fn func(albegas.Context), data ...interface{}) {
	cb := js.FuncOf(func(_this js.Value, args []js.Value) interface{} {
		var detail interface{} = data
		if len(data) == 1 {
			detail = data[0]
		}
		ctx := newContext(c, element.Of(args[0].Get("target")), args[0], detail)
		fn(ctx)
		return nil
	})

	c.on(elm, typ, cb)
}

// -----------------------------------------------------------------------------

func (c *component) Off(elm, typ string) {
	key := typKey(elm, typ)
	cb, ok := c.callbacks[key]
	if ok && cb.Truthy() {
		target := c.target(elm)
		target.RemoveEventListener(typ, cb)
		delete(c.callbacks, key)
	}
}

// -----------------------------------------------------------------------------

func (c *component) AddEventListener(typ string, listener js.Func) {
	c.on(Self, typ, listener)
}

// -----------------------------------------------------------------------------

func (c *component) RemoveEventListener(typ string, _ js.Func) {
	c.Off(Self, typ)
}

// -----------------------------------------------------------------------------

func (c *component) Depose() {
	for _, fn := range c.callbacks {
		fn.Release()
	}

	_app.remove(c.tatto)
}

// -----------------------------------------------------------------------------

func createComponentByElement(elm albegas.Element) *component {
	tatto := tatto(10)
	elm.SetAttr("data-albegas", tatto)
	ret := &component{
		tatto:     tatto,
		Element:   elm,
		callbacks: make(map[string]js.Func),
	}

	_app.append(tatto, ret)
	return ret
}

// -----------------------------------------------------------------------------

func createComponentByJSValue(value js.Value) *component {
	return createComponentByElement(element.Of(value))
}

// -----------------------------------------------------------------------------

func createComponetByTemplate(value js.Value) *component {

	if !value.Truthy() {
		panic(fmt.Sprintf("template is not valid"))
	}

	content := value.Get("content")

	return createComponentByJSValue(
		content.Call("cloneNode", true).Get("firstElementChild"),
	)
}

// -----------------------------------------------------------------------------

func CreateComponentByTemplate(id string) albegas.Component {
	return createComponetByTemplate(global.Query(id))
}

// -----------------------------------------------------------------------------

func CreateComponentByHTML(content albegas.HTML) albegas.Component {
	tmpl := global.CreateTemplate()
	tmpl.Set("innerHTML", content)
	return createComponetByTemplate(tmpl)
}

// -----------------------------------------------------------------------------

func CreateComponentById(id string) albegas.Component {
	return createComponentByElement(element.Of(global.Query(id)))
}
