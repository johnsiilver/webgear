package html

import (
	"html/template"
)

// Direction specifies the text direction for the content in an element.
type Direction string

const (
	// LTRDir is Left-to-right text direction.
	LTRDir = "ltr"
	// RTLDir Right-to-left text direction.
	RTLDir = "rtl"
	// AutoDir Let the browser figure out the text direction, based on the content (only recommended if the text direction is unknown).
	AutoDir = "auto"
)

// YesNo provides a yes or no for an attribute.
type YesNo string

func (y YesNo) String() string {
	if y == "no" {
		return ""
	}
	return "yes"
}

const (
	Yes YesNo = "yes"
	No  YesNo = "no"
)

// GlobalAttr is global attributes tha can be assigned to various elements.
type GlobalAttrs struct {
	// AccessKey specifies a shortcut key to activate/focus an element.
	AccessKey string
	// Class specifies one or more classnames for an element (refers to a class in a style sheet).
	Class string
	// ContentEditable specifies whether the content of an element is editable or not.
	ContentEditable bool
	// Dir specifies the text direction for the content in an element.
	Dir Direction
	// Draggable specifies whether the element is draggable or not.
	Draggable bool
	// Hidden specifies that an element is not yet, or is no longer, relevant.
	Hidden bool `html:"attr"`
	// ID specifies a unique id for an element.
	ID string
	// Lang specifies the language of the element's content.
	Lang string
	// SpellCheck specifies whether the element is to have its spelling and grammar checked or not.
	SpellCheck bool
	// Style specifies an inline CSS style for an element.
	Style string
	// TabIndex specifies the tabbing order of an element.
	// Note: we treat the index different than the standard. 0 will not be written out as a value. We required you
	// to have a sane order (it is not sane to have tab order 1, 2, 3, 4, 5, 0).
	TabIndex int
	// Title specifies extra information about an element.
	Title string
	// Translate specifies extra information about an element.
	Translate YesNo

	// XXXWasmUpdated indicates that a DocUpdater has updated the Element this is attached to, but
	// we have not flushed the changes to the DOM. This field is not a GlobalAttrs for HTML and is only public
	// so that is may be manipulated via reflection.
	XXXWasmUpdated bool
}

func (g GlobalAttrs) Attr() template.HTMLAttr {
	return template.HTMLAttr(structToString(g))
}
