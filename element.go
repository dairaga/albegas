//go:build js && wasm
// +build js,wasm

package albegas

import "syscall/js"

type HTML string

// -----------------------------------------------------------------------------

type ElementList []Element

// -----------------------------------------------------------------------------

type Element interface {
	Wrapper

	Prop(string) js.Value
	SetProp(string, interface{}) Element

	Attr(string) string
	SetAttr(string, string) Element

	Text() string
	SetText(string) Element

	HTML() HTML
	SetHTML(HTML) Element

	Truthy() bool

	Has(string) bool
	Add(clz ...string) Element
	Remove(clz ...string) Element
	Toggle(string) Element
	Replace(string, string) Element

	Value() string
	SetValue(string) Element

	QueryAll(string) ElementList
	Query(string) Element

	AddEventListener(string, js.Func)
	RemoveEventListener(string, js.Func)
	DispatchEvent(js.Value)
}

// -----------------------------------------------------------------------------
