/*
Package html provide objects representing HTML 5 tags that can be joined into a *Doc object and rendered.
It provides support for attaching event objects and dynamic content.

This package is akin to assembly language for web display.
*/
//
// Prerequisites
//
/*
To get the most out of the package, you will need the following:
	* An understanding of how HTML/CSS works
	* More than a passing familiarity with the net/http package
*/
//
// Package Layout
//
/*
All HTML tags and attribute types are contained in one large package. Normally this is undesirable,
creating too much "noise" in package documentation. But we want to emulate the
HTML flat namespace to allow easy browsing for people already familiar with HTML.

If you are looking
for an <a> tag, it will be html.A{}. If you are looking for an <img>, it will be html.Img{} and so on.

Note: This package does not currently contain all tags in HTML 5.  These will be slowly added.
*/
//
// Element - The Basic Type
//
/*
The Element interface is implemeted by all all html objects that wil produce output. These objects
represent our HTML tags and special object that can output dynamic content.
*/
//
// Doc - The Core Type
//
// The Doc type represents an HTML 5 document. Everything in the package eventually will end up in in a Doc for
// rendering.
//
// The Doc is also used to build component.Gear objects which represent HTML Shadow-DOM components.
//
// Basic Usage
//
// Create a page with basic header information that links to a stylesheet and prints "Hello World".
//	doc := &html.Doc{
// 		Head: &html.Head{
// 			Elements: []html.Element{
// 				&html.Meta{Charset: "UTF-8"},
// 				&html.Title{TagValue: html.TextElement("Simple Example")},
// 				&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/index.css")},
// 				&html.Link{Href: html.URLParse("https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap"), Rel: "stylesheet"},
// 			},
// 		},
// 		Body: &html.Body{
// 			Elements: []html.Element{
// 				html.TextElement("Hello World")
// 			},
// 		},
// 	}
//
// 	// To serve the page, see the handlers/ package.
//
// Serve Dynamic Content
//
/*
Similar to the template system, this package allows adding dynamic generated content. This is done by
implementing a DynamicFunc type and attaching it as an Element object to another type. This is done
via the Dynamic type.
A DynamicFunc simply creates and returns a list of []Element objects that will be inserted.
This example will serve a page that dynamically generates content from a slice. In a real situation, this could
be fetched from a file, a cache, ....
	// site represents a name of a site and a url to that site.
	type site struct {
		name string
		url *url.URL
	}

	// sites is used to implement html.DynamicFunc.
	type sites struct {
		sites []site
	}

	// List implements html.DynamicFunc to take a list of sites and display links to them.
	func (s sites) List(pipe Pipeline) []html.Element {
		if err := pipe.Ctx.Err(); err != nil {
			return nil
		}
		elements := []html.Element{}
		for _, site := range s.sites {
			elements = append(
				elements,
				&html.A{Href: url.String(), Elements: []html.Element(html.TextElement(site.name))},
			)
		}
		return elements
	}

	// This creats an instance of sites{} that has various names and urls that we wish to list out.
	var siteLister = sites{
	sites: []site{
		{"Microsoft", html.URLParse("www.microsoft.com")},
		{"Google", html.URLParse("www.google.com")},
		{"LucasFilm", html.URLParse("www.lucasfilm.com")},
	}

	// Create a page with basic header information that links to a stylesheet and prints "Hello World".
	doc := &html.Doc{
		Head: &html.Head{
			Elements: []html.Element{
				&html.Meta{Charset: "UTF-8"},
				&html.Title{TagValue: html.TextElement("Dyanic Example")},
				&html.Link{Rel: "stylesheet", Href: html.URLParse("/static/index.css")},
				&html.Link{Href: html.URLParse("https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap"), Rel: "stylesheet"},
			},
		},
		Body: &html.Body{
			Elements: []html.Element{
				html.Dynamic(siteLister.List),
			},
		},
	}
	// To serve the page, see the handlers/ package.
*/
//
// The Pipeline
//
/*
When implementing DynamicFunc above, the function signature requires a Pipeline object.
Pipeline provides a changing context object that passes certain information through the program flow.
This includes the actual Context of a call, the http.Request object and the output buffer W. The other
attributes are not useful in this context.

Here is an example of a DynamicFunc that returns the path of the URL it received for this call:
	func URLPath(pipe Pipeline) []html.Element {
		if err := pipe.Ctx.Err(); err != nil {
			return nil
		}

		return []Element{html.TextElement(pipe.Req.URL.Path)}
	}
*/
//
// GlobalAttr
//
// Most objects implement the GlobalAttr. This provides standard attributes such as "id" or "class". When
// available, this is implemented by embedding a GlobalAttr object in a tag type.
//
//	&html.A{
//		GlobalAttr: GlobalAttr{ID: "golangLink"},
//		Elements: []Element{TextElement("Golang.org")},
//		Href: "http://golang.org",
//	}
//
// Attaching Events
//
// Events are available on objects that support them via an Event attribute.  You can attach an event by
// simply calling a method on the object that you wish to attach a script to.
//
//	&html.Body{Events: &Events{}.OnLoad("DoSomething")}
package html

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

// insideWasm indicates that we are executing inside a WASM environment. This is set to true by
// event_wasm.go's init() function that is guarded by build tags.
var insideWasm = false

var docTmpl = template.Must(template.New("doc").Parse(strings.TrimSpace(`
{{- if not .Self.Component}}<!DOCTYPE html>{{end}}
{{if not .Self.Component}}<html>
	{{.Self.Head.Execute .}}
{{- end}}
	{{.Self.Body.Execute .}}
{{- if not .Self.Component}}</html>{{end}}
`)))

// Pipeline represents a template pipeline. The Self attribute is only usable internally, any other use is
// not supported. Component is used only internally by the component pacakge, any other use is not supported.
// Data is what the user wishes to pass in for their application.
type Pipeline struct {
	// Ctx is the context of the call chain. This should be set by NewPipeline().
	Ctx    context.Context
	cancel context.CancelFunc

	// errCh provides a channel that stores the first error encountered while trying to execute our internal call
	// chain.
	errCh chan error

	// Req is the http.Request object for this call.
	Req *http.Request

	// W is the output buffer.
	W io.Writer

	// Self represents the data structure of the object that is executing the template. This allows
	// a template to access attributes that represent a tag, such as A{} accessing Href for rendering.
	// A user should not set this, as it is automatically changed by the various Element implementations.
	Self interface{}

	// GearData provides a map of pipeline data keyed by gear name.
	// TODO(johnsiilver): Might want to have the component package provide its own Pipeline that this
	// pipeline is embedded in. Then GearData would only belong to that pipelin.  GearData has no
	// affect on anything in this package.
	GearData interface{}
}

// NewPipeline creates a new Pipeline object.
func NewPipeline(ctx context.Context, req *http.Request, w io.Writer) Pipeline {
	ctx, cancel := context.WithCancel(ctx)

	return Pipeline{
		Ctx:    ctx,
		cancel: cancel,
		errCh:  make(chan error, 1),
		Req:    req,
		W:      w,
	}
}

// Error adds an error to the Pipeline. If there is already an error recorded, the error will be dropped.
func (p Pipeline) Error(err error) {
	p.cancel()
	select {
	case p.errCh <- err:
	default:
	}
}

// HadError returns an error if the pipeline had an error during execution.
func (p Pipeline) HadError() error {
	select {
	case err := <-p.errCh:
		return err
	default:
		return nil
	}
}

// Doc represents an HTML 5 document.
type Doc struct {
	Head *Head
	Body *Body

	// Pretty says to make the HTML look pretty before outputting.
	Pretty bool

	// Componenet is used to indicate that this is a snippet of code, not a full document.
	// As such, <html> and <head> tags will be suppressed.
	Component bool

	pool sync.Pool

	initDone bool
}

// Init sets up all the internals for execution. Must be called before Execute() and should only be called once.
func (d *Doc) Init() error {
	if d.initDone {
		return nil
	}

	if err := d.validate(); err != nil {
		return err
	}

	if d.Head != nil {
		if err := d.Head.Init(); err != nil {
			return err
		}
	}

	if d.Component {
		if d.Body != nil {
			d.Body.Component = true
		}
	}
	if err := d.Body.Init(); err != nil {
		return err
	}

	d.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	d.initDone = true
	return nil
}

// validate attempts to do basic validation of the Doc contents as best it can.
func (d *Doc) validate() error {
	if err := d.Body.validate(); err != nil {
		return err
	}
	return nil
}

// Element represents an object that can render self container HTML 5. Normally this is an HTML5 tag.
// Users may implement this, but do so at their own risk as we can change the implementation without
// changing the major version.
type Element interface {
	// Execute outputs the Element's textual representation to Pipeline.W . Execute returns a string,
	// but that string value MUST always be an empty string. This is a side effect of the Go template
	// system not allowing a function call unless it provides output or is in a FuncMap. FuncMap
	// is not usable in this context.
	Execute(pipe Pipeline) string
}

// EmptyString is returned by all Element.Execute() calls.
const EmptyString = ""

// Initer is a type that requires Init() to be called before using.
type Initer interface {
	// Init initalizes the internal state.
	Init() error
}

// outputAble details if a slice or struct should be output to when doing structToString.
type outputAble interface {
	outputAble()
	fmt.Stringer
}

// raw details if a type should have its output from String() just shoved in without having the stuct's
// field name rendered or anything else.
type raw interface {
	isRaw()
	fmt.Stringer
}

// DynamicFunc is a function that uses dynamic server data to return Elements that will be rendered.
type DynamicFunc func(pipe Pipeline) []Element

type dynamic struct {
	f DynamicFunc
}

func (d *dynamic) Execute(pipe Pipeline) string {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("DynamicFunc with type %T paniced: %s\nstack trace:\n%s", d.f, r, string(debug.Stack()))
		}
	}()

	pipe.Self = d

	elements := d.f(pipe)
	compileElements(elements)
	for _, e := range elements {
		if pipe.Ctx.Err() != nil {
			return EmptyString
		}
		e.Execute(pipe)
	}
	return EmptyString
}

// Dynamic wraps a DynamicFunc so that it implements Element.
func Dynamic(f DynamicFunc) Element {
	return &dynamic{
		f: f,
	}
}

// TextElement is an element that represents text, usually in a value. It is not valid everywhere.
type TextElement string

func (t TextElement) Execute(pipe Pipeline) string {
	pipe.W.Write([]byte(t))
	return EmptyString
}

func (t TextElement) isZero() bool {
	if t == "" {
		return true
	}
	return false
}

func structToString(i interface{}) string {
	val := reflect.ValueOf(i)

	// If it is *struct, get the struct and assign back to val.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic(fmt.Sprintf("structToString() received %T instead of a struct or *struct", i))
	}

	out := []string{}

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		isNaked := false // Details if this is supposedly a attribute with no value attached. Such as <tag attr> vs <tag attr=true>.

		// Retrieve the field name.
		name := ""
		sf := t.Field(i)
		if sf.Anonymous {
			continue
		}

		// While we handle skipping Elements when their type is an interface or
		// some slice like []Element, those don't handle when we do specific types
		// like TextElement (which is based on string).  So this part eliminates
		// key words that we know shouldn't get rendered as an attribute.
		switch sf.Name {
		case "Element", "Elements":
			continue
		}

		// Names starting with _ are for internal use.
		if strings.HasPrefix(sf.Name, "XXX") {
			continue
		}

		// Non-exported field.
		if strings.Title(sf.Name) != sf.Name {
			continue
		}

		if tagName := sf.Tag.Get("html"); tagName != "" {
			// Special: it says that it is a tag without value, like <tag attribute>.
			if tagName == "attr" {
				isNaked = true
				name = sf.Name
			} else {
				name = tagName
			}
		} else {
			name = sf.Name
		}

		// Special value that we skip.
		// TODO(johnsiilver): Not sure we still need this.
		if name == "TagValue" {
			continue
		}

		field := val.Field(i)

		// This handles the case where we just want to put the raw output of the String() method.
		if r, ok := field.Interface().(raw); ok {
			o := r.String()
			if o == "" {
				continue
			}
			out = append(out, o)
		}

		// Retrieve the value.
		var str string
		switch field.Kind() {
		case reflect.String:
			str = field.String()
			// This detects that the field was a zero value string, so don't include it.
			if str == "" {
				continue
			}
		case reflect.Bool:
			switch field.Bool() {
			case true:
				str = "true"
			case false:
				// We don't do zero values.
				continue
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() == 0 {
				continue
			}
			str = strconv.FormatInt(field.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if field.Uint() == 0 {
				continue
			}
			str = strconv.FormatUint(field.Uint(), 10)
		case reflect.Struct, reflect.Slice:
			_, ok := field.Interface().(outputAble)
			if !ok {
				continue
			}
			meth := field.MethodByName("String")
			str = meth.Call(nil)[0].Interface().(string)
			if str == "" {
				continue
			}
		case reflect.Ptr:
			if field.IsNil() {
				continue
			}
			v := field.Elem()

			if fv, ok := v.Interface().(url.URL); ok {
				str = fv.String()
			} else {
				continue
			}
		case reflect.Interface:
			continue
		default:
			if meth := field.MethodByName("Get" + t.Name()); meth.IsValid() {
				str = meth.Call(nil)[0].Interface().(string)
			} else if meth := field.MethodByName("String" + t.Name()); meth.IsValid() {
				str = meth.Call(nil)[0].Interface().(string)
			} else {
				panic(fmt.Sprintf("structToString on field %q: non-string , int, uint or getter method, is %s", name, field.Kind()))
			}
		}

		// Handles when we want to add something like "px" or "em" to the end of a number.
		if suffix := sf.Tag.Get("suffix"); suffix != "" {
			str = str + suffix
		}

		// Naked is about if the attribute should be just a tag with no "=", like "sandbox" instead of "sandbox=".
		if isNaked {
			out = append(out, strings.ToLower(name))
		} else {
			out = append(out, fmt.Sprintf("%s=%q", strings.ToLower(name), str))
		}
	}

	return strings.Join(out, " ")
}

// compileElements compiles every Element passed and recursively all Elements contained in those Element(s).
func compileElements(elements []Element) error {
	for _, element := range elements {
		if err := compileElement(element); err != nil {
			return err
		}
	}
	return nil
}

// compileElement complies the passed Element and recursively all Elements contained in that Element.
func compileElement(element Element) error {
	if i, ok := element.(Initer); ok {
		if err := i.Init(); err != nil {
			return err
		}
	}

	val := reflect.ValueOf(element)

	// If it is *struct, get the struct and assign back to val.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := val.Field(i)

		if t.Field(i).Anonymous || !field.CanInterface() {
			continue
		}

		switch real := field.Interface().(type) {
		case Element:
			if err := compileElement(real); err != nil {
				return nil
			}
		case []Element:
			if err := compileElements(real); err != nil {
				return nil
			}
		}
	}
	return nil
}

// URLParse returns a *url.URL representation of "s". If it cannot be parsed, this will panic.
func URLParse(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return u
}
