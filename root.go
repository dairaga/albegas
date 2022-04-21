//go:build js && wasm
// +build js,wasm

package albegas

import "syscall/js"

type Wrapper interface {
	Ref() js.Value
}
