//go:build js && wasm
// +build js,wasm

package element

import (
	"syscall/js"

	"github.com/dairaga/albegas"
)

type element struct {
	ref js.Value
}

// -----------------------------------------------------------------------------

func (e *element) Ref() js.Value {
	return e.ref
}

// -----------------------------------------------------------------------------

func (e *element) Prop(p string) js.Value {
	return e.ref.Get(p)
}

// -----------------------------------------------------------------------------

func (e *element) SetProp(p string, v interface{}) albegas.Element {
	e.ref.Set(p, v)
	return e
}

// -----------------------------------------------------------------------------

func (e *element) Attr(name string) string {
	value := e.ref.Call("getAttribute", name)
	if value.Truthy() {
		return value.String()
	}
	return ""
}

// -----------------------------------------------------------------------------

func (e *element) SetAttr(name, value string) albegas.Element {
	e.ref.Call("setAttribute", name, value)
	return e
}

// -----------------------------------------------------------------------------

func (e *element) Text() string {
	value := e.Prop("innerText")
	if value.Truthy() {
		return value.String()
	}
	return ""
}

// -----------------------------------------------------------------------------

func (e *element) SetText(text string) albegas.Element {
	return e.SetProp("innerText", text)
}

// -----------------------------------------------------------------------------

func (e *element) HTML() albegas.HTML {
	value := e.Prop("innerHTML")
	if value.Truthy() {
		return albegas.HTML(value.String())
	}
	return albegas.HTML("")
}

// -----------------------------------------------------------------------------

func (e *element) SetHTML(content albegas.HTML) albegas.Element {
	e.SetProp("innerHTML", content)
	return e
}

// -----------------------------------------------------------------------------

func (e *element) Truthy() bool {
	return e.ref.Truthy()
}

// -----------------------------------------------------------------------------

func (e *element) clz(method string, args ...string) albegas.Element {
	size := len(args)
	if size <= 0 {
		return e
	}

	if size == 1 {
		e.ref.Get("classList").Call(method, args[0])
	} else if size == 2 {
		e.ref.Get("classList").Call(method, args[0], args[1])
	} else if size == 3 {
		e.ref.Get("classList").Call(method, args[0], args[1], args[2])
	} else {
		x := make([]interface{}, size)
		for i, str := range args {
			x[i] = str
		}
		e.ref.Get("classList").Call(method, x...)
	}

	return e
}

// -----------------------------------------------------------------------------

// AddClass add class to element.
func (e *element) Add(names ...string) albegas.Element {
	return e.clz("add", names...)
}

// -----------------------------------------------------------------------------

// RemoveClass remove class from element.
func (e *element) Remove(names ...string) albegas.Element {
	return e.clz("remove", names...)
}

// -----------------------------------------------------------------------------

// ToggleClass toggle some class of element.
func (e *element) Toggle(name string) albegas.Element {
	return e.clz("toggle", name)
}

// -----------------------------------------------------------------------------

// ReplaceClass replace some class of element with new one.
func (e *element) Replace(oldName, newName string) albegas.Element {
	return e.clz("replace", oldName, newName)
}

// -----------------------------------------------------------------------------

// HasClass returns boolean indicates whether or not element has the class.
func (e *element) Has(name string) bool {
	value := e.Prop("classList")
	if value.Truthy() {
		return false
	}
	return value.Call("contains", name).Bool()
}

// -----------------------------------------------------------------------------

func (e *element) Value() string {
	value := e.Prop("value")
	if value.Truthy() {
		return ""
	}

	return value.String()
}

// -----------------------------------------------------------------------------

func (e *element) SetValue(value string) albegas.Element {
	e.SetProp("value", value)
	return e
}

// -----------------------------------------------------------------------------

func (e *element) QueryAll(selector string) albegas.ElementList {
	value := e.ref.Call("querySelectorAll", selector)
	size := value.Length()
	if size <= 0 {
		return albegas.ElementList([]albegas.Element{})
	}

	lst := make([]albegas.Element, size)
	for i := 0; i < size; i++ {
		lst[i] = Of(value.Index(i))
	}
	return lst
}

// -----------------------------------------------------------------------------

func (e *element) Query(selector string) albegas.Element {
	value := e.ref.Call("querySelector", selector)
	if value.Truthy() {
		return Of(value)
	}
	return nil
}

// -----------------------------------------------------------------------------

func (e *element) AddEventListener(typ string, listener js.Func) {
	e.ref.Call("addEventListener", typ, listener)
}

// -----------------------------------------------------------------------------

func (e *element) RemoveEventListener(typ string, listener js.Func) {
	e.ref.Call("removeEventListener", typ, listener)
}

// -----------------------------------------------------------------------------

func (e *element) DispatchEvent(event js.Value) {
	e.ref.Call("dispatchEvent", event)
}

// -----------------------------------------------------------------------------

func Of(val js.Value) albegas.Element {
	return &element{
		ref: val,
	}
}
