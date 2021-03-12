// +build js,wasm

package html

import (
	html "html/template"
)

// GearType is an interface only meant to be implemented by *component.Gear. We cannot
// use component.Gear because component uses this package (cyclic dependency). We do
// not want to merge these packages and don't want to migrate the Element type out.
// So this is used as a stand in. The only valid use for this in client code is for tests,
// which should always embed GearType.  Any other use has no compatibility promise.
type GearType interface {
	Element
	Name() string
	GearID() string
	TagType() html.HTMLAttr
	Execute(pipe Pipeline) string
	UpdateFlag() bool
	SetUpdateFlag()
	RemoveUpdateFlag()
	UpdateDOM() error
}
