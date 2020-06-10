package html

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/yosssi/gohtml"
)

var docTmpl = template.Must(template.New("doc").Parse(strings.TrimSpace(`
{{- if not .Self.Component}}<html>
	{{.Self.Head.Execute .}}
{{- end}}
	{{.Self.Body.Execute .}}
{{- if not .Self.Component}}</html>{{end}}
`)))

// Pipeline represents a template pipeline. The Self attribute is only usable internally, any other use is
// not supported. Component is used only internally by the component pacakge, any other use is not supported.
// Data is what the user wishes to pass in for their application.
type Pipeline struct {
	// Req is the http.Request object for this call.
	Req *http.Request
	// W is the http.ResponseWriter. This can be used with the http.Error() call in the standard lib to
	// indicate an error.
	W http.ResponseWriter
	// Self represents the data structure of the object that is executing the template. This allows
	// a template to access attributes that represent a tag, such as A{} accessing Href for rendering.
	// A user should not set this, as it is automatically changed by the various Element implementations.
	Self interface{}
	// GearData provides a map of pipeline data keyed by gear name.
	GearData map[string]interface{}
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
}

// Init sets up all the internals for execution.  Must be called before Execute(). You only need to do this once.
func (d *Doc) Init() error {
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

	return nil
}

// Execute executes the internal template and writes the output to the io.Writer.
func (d *Doc) Execute(w io.Writer, pipe Pipeline) error {
	pipe.Self = d

	if d.Pretty {
		return docTmpl.ExecuteTemplate(gohtml.NewWriter(w), "doc", pipe)
	}

	return docTmpl.ExecuteTemplate(w, "doc", pipe)
}

// Render calls execute execute and returns the string value.
func (d *Doc) Render(pipe Pipeline) (template.HTML, error) {
	w := d.pool.Get().(*strings.Builder)
	defer d.pool.Put(w)
	w.Reset()

	if err := d.Execute(w, pipe); err != nil {
		return "", err
	}

	return template.HTML(w.String()), nil
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
	Execute(pipe Pipeline) template.HTML
}

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

// raw details if a type should have its output from String() just shoved in without field names or anything.
type raw interface {
	isRaw()
	fmt.Stringer
}

// DynamicFunc is a function that uses dynamic server data to return Elements that will be rendered.
type DynamicFunc func(pipe Pipeline) []Element

type dynamic struct {
	f    DynamicFunc
	pool sync.Pool
}

func (d *dynamic) Execute(pipe Pipeline) template.HTML {
	buff := d.pool.Get().(*strings.Builder)
	defer d.pool.Put(buff)
	buff.Reset()

	pipe.Self = d

	elements := d.f(pipe)
	compileElements(elements)
	for _, e := range elements {
		buff.WriteString(string(e.Execute(pipe)))
	}
	return template.HTML(buff.String())
}

// Dynamic wraps a DynamicFunc so that it implements Element.
func Dynamic(f DynamicFunc) Element {
	return &dynamic{
		f: f,
		pool: sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}
}

// TextElement is an element that represents text, usually in a value. It is not valid everywhere.
type TextElement string

func (t TextElement) Execute(pipe Pipeline) template.HTML {
	return template.HTML(t)
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
		isNaked := false // Details if this is supposedly a attribute with no value attached.

		// Retrieve the field name.
		name := ""
		sf := t.Field(i)
		if sf.Anonymous {
			continue
		}

		// Non-exported field.
		if strings.Title(sf.Name) != sf.Name {
			continue
		}

		if tagName := sf.Tag.Get("html"); tagName != "" {
			// Special: it says that it is a tag without value.
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
			// This detects that the field was a zero value, so don't include it.
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

// compileElement complies the passed Element and ecursively all Elements contained in that Element.
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
