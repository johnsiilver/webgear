package html

import (
	"fmt"
)

// FormMethod is the method the form will use to submit the form's content.
type FormMethod string

const (
	// GetMethod uses the browser's "get" method.
	GetMethod FormMethod = "get"
	// PostMethod uses the browser's "post" method.
	PostMethod FormMethod = "post"
)

// FormElement represents an element within an HTML form.
type FormElement interface {
	isFormElement()
}

// Form represents an HTML form.
type Form struct {
	// Action defines the action to be performed when the form is submitted.
	// If the action attribute is omitted, the action is set to the current page.
	Action string

	// Target attribute specifies if the submitted result will open in a new browser tab, a frame,
	// or in the current window. The default value is "_self" which means the form will be submitted
	// in the current window. To make the form result open in a new browser tab, use the value "_blank".
	// Other legal values are "_parent", "_top", or a name representing the name of an iframe.
	Target string

	// Method attribute specifies the HTTP method (GET or POST) to be used when submitting the form data.
	Method FormMethod

	// AcceptCharset specifies the charset used in the submitted form (default: the page charset).
	AcceptCharset string `html:"accept-charset"`

	// AutoComplete specifies if the browser should autocomplete the form (default: on).
	AutoComplete bool

	// Enctype pecifies the encoding of the submitted data (default: is url-encoded).
	EncType string

	// NoValidate Specifies that the browser should not validate the form.
	NoValidate bool

	// FormElements are elements contained form.
	FormElements []FormElement

	events Events
}

// Events returns the internal events object that allows you to attached events to trigger Javascript functions.
func (f *Form) Events() *Events {
	return &f.events
}

func (f *Form) isElement() {}

// InputType describes the type of input that is being created within a form.
type InputType string

const (
	// TextInput takes in text from a keyboard.
	TextInput InputType = "text"
	// RadioInput creates a radio button.
	RadioInput InputType = "radio"
	// SubmitInput creates a submit button.
	SubmitInput InputType = "submit"
)

// Input creates a method of input within a form.
type Input struct {
	GlobalAttrs

	// Type is the type of input.
	Type InputType
	// Name is the name of the input field, this is mandatory.
	Name string
	// Value is the value of the input field.
	Value string
}

func (i Input) validate() error {
	if i.Name == "" {
		return fmt.Errorf("a Form had an Input without .Name, which is an error")
	}
	return nil
}

func (i Input) isFormElement() {}

// Label element is useful for screen-reader users, because the screen-reader will read out loud the label
// when the user is focused on the input element.
type Label struct {
	For   string
	Value string
}

func (l Label) isFormElement() {}
