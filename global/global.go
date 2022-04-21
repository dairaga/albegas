//go:build js && wasm
// +build js,wasm

package global

import (
	"syscall/js"
)

var _global = js.Global()
var _window = _global.Get("window")
var _document = _window.Get("document")
var _body = _document.Get("body")

// -----------------------------------------------------------------------------

// CreateElement returns a HTML element.
func CreateElement(tag string) js.Value {
	return _document.Call("createElement", tag)
}

// -----------------------------------------------------------------------------

func CreateTemplate() js.Value {
	return CreateElement("template")
}

// -----------------------------------------------------------------------------

// AppendChild appends child to document.
func AppendChild(child js.Value) js.Value {
	return _body.Call("appendChild", child)
}

// -----------------------------------------------------------------------------

func Query(selector string) js.Value {
	return _document.Call("querySelector", selector)
}

// -----------------------------------------------------------------------------

func QueryAll(selector string) js.Value {
	return _document.Call("querySelectorAll", selector)
}
