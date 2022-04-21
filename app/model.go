//go:build js && wasm
// +build js,wasm

package app

type model struct {
	name  string
	value interface{}
}

// -----------------------------------------------------------------------------

func (m *model) Name() string {
	return m.name
}

// -----------------------------------------------------------------------------

func (m *model) Value() interface{} {
	return m.value
}
