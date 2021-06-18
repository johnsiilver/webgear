package html

import (
	"fmt"
	"html/template"
	"strings"
)

var textareaTmpl = template.Must(template.New("textarea").Parse(strings.TrimSpace(`
<textarea {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}} {{.Self.Attr}}>
	{{- $data := .}}
	{{- .Self.Element.Execute $data}}</textarea>
`)))

// WrapType specifies how the text in a textarea is to be wrapped when submitted in a form.
type WrapType string

const (
	// SoftWrap specifies the text in the textarea is not wrapped when submitted in a form. This is default.
	SoftWrap WrapType = "soft"
	// HardWrap specifies the text in the textarea is wrapped (contains newlines)
	// when submitted in a form. When "hard" is used, the cols attribute must be specified.
	HardWrap WrapType = "hard"
)

// TextArea tag defines a multi-line text input control.
type TextArea struct {
	GlobalAttrs
	Events *Events

	// Name specifies a name for a text area.
	Name string
	// Form specifies which form the text area belongs to.
	Form string

	// Cols specifies the visible width of a text area.
	Cols int
	// MaxLength specifies the maximum number of characters allowed in the text area.
	MaxLength int
	// Rows specifies the visible number of lines in a text area.
	Rows int

	// DirName specifies that the text direction of the textarea will be submitted.
	DirName string
	// Wrap specifies how the text in a text area is to be wrapped when submitted in a form.
	Wrap WrapType
	// Placeholder specifies a short hint that describes the expected value of a text area.
	Placeholder string

	// AutoFocus specifies that a text area should automatically get focus when the page loads.
	AutoFocus bool `html:"attr"`
	// Disabled specifies that a text area should be disabled.
	Disabled bool `html:"attr"`
	// ReadOnly specifies that a text area should be read-only.
	ReadOnly bool `html:"attr"`
	// Required specifies that a text area is required/must be filled out.
	Required bool `html:"attr"`

	// Element contains a TextElement that is the content of TextArea.
	Element TextElement
}

func (t *TextArea) validate() error {
	if t.DirName != "" {
		switch {
		case t.Name == "":
			return fmt.Errorf("TextArea.DirName(%s) is set, but .Name is not", t.DirName)
		case t.DirName != t.Name+".dir":
			return fmt.Errorf("TextArea.DirName(%s) is set and not == to %s", t.DirName, t.Name+".dir")
		}
	}
	return nil
}

func (t *TextArea) Attr() template.HTMLAttr {
	output := structToString(t)
	return template.HTMLAttr(output)
}

func (t *TextArea) Execute(pipe Pipeline) string {
	pipe.Self = t

	if err := textareaTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
