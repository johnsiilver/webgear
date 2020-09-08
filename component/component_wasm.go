// +build js,wasm

package component

import (
	"syscall/js"
)

// UpdateDOM updates the DOM for this component.
func (g *Gear) UpdateDOM() error {
	js.Global().Get("document").Call("getElementById", g.TemplateName()).Set("outerHTML", g.TemplateContent())
	js.Global().Call(string(g.LoaderName()))
	return nil
}

// UpdateFlag returns the gears wasm update flag. This is not meant for users.
func (g *Gear) UpdateFlag() bool {
	g.wasmUpdateMu.Lock()
	defer g.wasmUpdateMu.Unlock()
	return g.wasmUpdate
}

// SetUpdateFlag sets the gears wasm update flag to true. This is not meant for users.
func (g *Gear) SetUpdateFlag() {
	g.wasmUpdateMu.Lock()
	defer g.wasmUpdateMu.Unlock()
	g.wasmUpdate = true
}

// RemoveUpdateFlag removes the gear's wasm update flag to false. This is not mean for users.
func (g *Gear) RemoveUpdateFlag() {
	g.wasmUpdateMu.Lock()
	defer g.wasmUpdateMu.Unlock()
	g.wasmUpdate = false
}
