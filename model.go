//go:build js && wasm
// +build js,wasm

package albegas

type Model interface {
	Name() string
	Value() interface{}
}
