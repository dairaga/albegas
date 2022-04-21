//go:build js && wasm
// +build js,wasm

package albegas

type Component interface {
	Element
	Append(string, Component)
	On(string, string, func(Context), ...interface{})
	Off(string, string)
	Watch(string, interface{})
	WatchIndex(string, int, interface{})
	WatchMapIndex(string, interface{}, interface{})
	Depose()
}
