//go:build js && wasm
// +build js,wasm

package albegas

import "syscall/js"

type Context interface {
	Owner() Component
	Sender() Element
	Event() js.Value
	Detail() interface{}
}
