package html

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

var formTmpl = template.Must(template.New("form").Parse(strings.TrimSpace(`
<form {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</form>
`)))

// FormMethod is the method the form will use to submit the form's content.
type FormMethod string

const (
	// GetMethod uses the browser's "get" method.
	GetMethod FormMethod = "get"
	// PostMethod uses the browser's "post" method.
	PostMethod FormMethod = "post"
)

// FormRel specifies the relationship between a linked resource and the current document in a form.
type FormRel string

const (
	ExternalFR   FormRel = "external"
	HelpFR       FormRel = "help"
	LicenseFR    FormRel = "license"
	NextFR       FormRel = "next"
	NofollowFR   FormRel = "nofollow"
	NoopenerFR   FormRel = "noopener"
	NoreferrerFR FormRel = "noreferrer"
	OpenerFR     FormRel = "opener"
	PrevFR       FormRel = "prev"
	SearchFR     FormRel = "search"
)

// FormElement represents an element within an HTML form.
type FormElement interface {
	Element
	isFormElement()
}

// Form represents an HTML form.
type Form struct {
	GlobalAttrs

	Events *Events

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
	NoValidate bool `html:"attr"`

	// Rel specifies the relationship between a linked resource and the current document.
	Rel FormRel

	// Elements are elements contained form.
	Elements []FormElement
}

func (f *Form) Execute(pipe Pipeline) string {
	pipe.Self = f

	if err := formTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (f *Form) Attr() template.HTMLAttr {
	output := structToString(f)
	return template.HTMLAttr(output)
}

var inputTmpl = template.Must(template.New("input").Parse(strings.TrimSpace(`
<input {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
`)))

// InputType describes the type of input that is being created within a form.
type InputType string

const (
	// TextInput takes in text from a keyboard.
	TextInput InputType = "text"
	// RadioInput creates a radio button.
	RadioInput InputType = "radio"
	// SubmitInput creates a submit button.
	SubmitInput InputType = "submit"
	// ButtonInput creates a button.
	ButtonInput InputType = "button"
)

// Input creates a method of input within a form.
type Input struct {
	GlobalAttrs
	*Events

	// Type is the type of input.
	Type InputType
	// Name is the name of the input field, this is mandatory.
	Name string
	// Value is the value of the input field.
	Value string
}

func (i Input) isFormElement() {}

func (i *Input) Execute(pipe Pipeline) string {
	pipe.Self = i

	if err := inputTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (i *Input) Attr() template.HTMLAttr {
	output := structToString(i)
	return template.HTMLAttr(output)
}

func (i Input) validate() error {
	if i.Name == "" {
		return fmt.Errorf("a Form had an Input without .Name, which is an error")
	}
	return nil
}

var labelTmpl = template.Must(template.New("label").Parse(strings.TrimSpace(`
<label {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</label>
`)))

// Label element is useful for screen-reader users, because the screen-reader will read out loud the label
// when the user is focused on the input element.
type Label struct {
	GlobalAttrs
	*Events

	// For specifies the id of the form element the label should be bound to.
	For string
	// Form specifies which form the label belongs to.
	Form string

	// Elements are HTML elements that are contained in the Label tag.
	// Usually a TextElement and sometimes the input tag the Label is for.
	Elements []Element
}

func (l *Label) isFormElement() {}

func (l *Label) Execute(pipe Pipeline) string {
	pipe.Self = l

	if err := labelTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (l *Label) Attr() template.HTMLAttr {
	output := structToString(l)
	return template.HTMLAttr(output)
}

var buttonTmpl = template.Must(template.New("button").Parse(strings.TrimSpace(`
<button {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</button>
`)))

// ButtonType describes the kind of button a button tag is.
type ButtonType string

const (
	ButtonBT ButtonType = "button"
	ResetBT  ButtonType = "reset"
	SubmitBt ButtonType = "submit"
)

// Button implements the HTML button tag.
type Button struct {
	GlobalAttrs
	Events *Events

	// AutoFocus specifies that a text area should automatically get focus when the page loads.
	AutoFocus bool `html:"attr"`
	// Disabled specifies that a text area should be disabled.
	Disabled bool `html:"attr"`
	// Form specifies which form the text area belongs to.
	Form string
	// FormAction specifies where to send the form-data when a form is submitted. Only for type="submit".
	FormAction *url.URL
	// FormEncType secifies how form-data should be encoded before sending it to a server. Only for type="submit".
	FormEncType string
	// FormMethod specifies how to send the form-data (which HTTP method to use). Only for type="submit".
	FormMethod FormMethod
	// FormNoValidate specifies that the form-data should not be validated on submission. Only for type="submit".
	FormNoValidate bool `html:"attr"`
	// FormTarget specifies where to display the response after submitting the form. Only for type="submit".
	// Specific constants such as BlankTarget are defined for common target names in this package.
	FormTarget string
	// FrameName specifies where to display the response after submitting the form. Only for type="submit".
	FrameName string
	// Name specifies a name for the button.
	Name string
	// Type specifies the type of button.
	Type ButtonType
	// Value specifies an initial value for the button.
	Value string

	Elements []Element
}

func (b *Button) isFormElement() {}

func (b *Button) Execute(pipe Pipeline) string {
	pipe.Self = b

	if err := buttonTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (b *Button) Attr() template.HTMLAttr {
	output := structToString(b)
	return template.HTMLAttr(output)
}

var selectTmpl = template.Must(template.New("select").Parse(strings.TrimSpace(`
<select {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</select>
`)))

type SelectElement interface {
	Element
	isSelectElement()
}

type Select struct {
	GlobalAttrs
	Events *Events

	// AutoFocus specifies that a text area should automatically get focus when the page loads.
	AutoFocus bool `html:"attr"`
	// Disabled specifies that a text area should be disabled.
	Disabled bool `html:"attr"`
	// Form specifies which form the text area belongs to.
	Form string
	// Multiple specifies that multiple options can be selected at once.
	Multiple bool `html:"attr"`
	// Name specifies a name for the button.
	Name string
	// Required specifies that the user is required to select a value before submitting the form.
	Required bool `html:"attr"`
	// Size defines the number of visible options in a drop-down list.
	Size uint

	Elements []SelectElement
}

func (s *Select) isFormElement() {}

func (s *Select) Execute(pipe Pipeline) string {
	pipe.Self = s

	if err := selectTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (s *Select) Attr() template.HTMLAttr {
	output := structToString(s)
	return template.HTMLAttr(output)
}

var optionTmpl = template.Must(template.New("option").Parse(strings.TrimSpace(`
<option {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>{{.Self.TagValue}}</option>
`)))

// Option defines an option in a select list.
type Option struct {
	GlobalAttrs
	Events *Events

	// Disabled specifies that a text area should be disabled.
	Disabled bool `html:"attr"`
	// Lablel specifies a shorter label for an option.
	Label string
	// Selected specifies that an option should be pre-selected when the page loads.
	Selected bool `html:"attr"`
	// Value specifies the value to be sent to a server.
	Value string
	// TagValue is the text in the option.
	TagValue string
}

func (o *Option) isSelectElement() {}

func (o *Option) isOptGroupElement() {}

func (o *Option) Execute(pipe Pipeline) string {
	pipe.Self = o

	if err := optionTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (o *Option) Attr() template.HTMLAttr {
	output := structToString(o)
	return template.HTMLAttr(output)
}

var optGroupTmpl = template.Must(template.New("optgroup").Parse(strings.TrimSpace(`
<optgroup {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</optgroup>
`)))

type OptGroupElement interface {
	Element
	isOptGroupElement()
}

// OptGroup is used to group related options in a Select element.
type OptGroup struct {
	// Disabled specifies that a text area should be disabled.
	Disabled bool `html:"attr"`
	// Lablel specifies a shorter label for an option.
	Label string

	Elements []OptGroupElement
}

func (o *OptGroup) isSelectElement() {}

func (o *OptGroup) Execute(pipe Pipeline) string {
	pipe.Self = o

	if err := optGroupTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (o *OptGroup) Attr() template.HTMLAttr {
	output := structToString(o)
	return template.HTMLAttr(output)
}

var legendTmpl = template.Must(template.New("legend").Parse(strings.TrimSpace(`
<legend {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>{{.Self.Caption}}</legend>
`)))

// Legend defines a caption for the FieldSet element.
type Legend struct {
	// Caption is the legend's caption.
	Caption string
}

func (l *Legend) Execute(pipe Pipeline) string {
	pipe.Self = l

	if err := legendTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

var fieldSetTmpl = template.Must(template.New("fieldset").Parse(strings.TrimSpace(`
<fieldset {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</fieldset>
`)))

// FieldSet is used to group related elements in a form.
type FieldSet struct {
	// Disabled specifies that a text area should be disabled.
	Disabled bool `html:"attr"`
	// Form specifies which form the text area belongs to.
	Form string
	// Name specifies a name for the button.
	Name string

	Elements []Element
}

func (f *FieldSet) Execute(pipe Pipeline) string {
	pipe.Self = f

	if err := fieldSetTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (f *FieldSet) Attr() template.HTMLAttr {
	output := structToString(f)
	return template.HTMLAttr(output)
}

var outputTmpl = template.Must(template.New("output").Parse(strings.TrimSpace(`
<output {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}></output>
`)))

// Output is used to represent the result of a calculation.
type Output struct {
	// For specifies the relationship between the result of the calculation, and the elements used in the calculation.
	For string
	// Form specifies which form the text area belongs to.
	Form string
	// Name specifies a name for the button.
	Name string
}

func (o *Output) isFormElement() {}

func (o *Output) Execute(pipe Pipeline) string {
	pipe.Self = o

	if err := outputTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (o *Output) Attr() template.HTMLAttr {
	output := structToString(o)
	return template.HTMLAttr(output)
}
