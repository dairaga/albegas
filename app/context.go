//go:build js && wasm
// +build js,wasm

package app

import (
	"syscall/js"

	"github.com/dairaga/albegas"
)

type context struct {
	owner  albegas.Component
	sender albegas.Element
	event  js.Value
	detail interface{}
}

// -----------------------------------------------------------------------------

func (c *context) Owner() albegas.Component {
	return c.owner
}

// -----------------------------------------------------------------------------

func (c *context) Sender() albegas.Element {
	return c.sender
}

// -----------------------------------------------------------------------------

func (c *context) Event() js.Value {
	return c.event
}

// -----------------------------------------------------------------------------

func (c *context) Detail() interface{} {
	return c.detail
}

// -----------------------------------------------------------------------------

func newContext(owner *component, sender albegas.Element, event js.Value, detail interface{}) *context {
	return &context{
		owner:  owner,
		sender: sender,
		event:  event,
		detail: detail,
	}
}
