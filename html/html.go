package html

import (
	"fmt"
	"html/template"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/yosssi/gohtml"
)

var docTmpl = strings.TrimSpace(`
{{- if not .Self.Component}}<html>
	{{.Self.Head.Execute .Data }}
{{- end}}
	{{.Self.Body.Execute .Data}}
{{- if not .Self.Component}}</html>{{end}}
`)

type pipeline struct {
	Self interface{}
	Data interface{}
}

// Doc represents an HTML 5 document.
type Doc struct {
	Head *Head
	Body *Body

	// Pipeline provides a Go template pipeline object.
	Pipeline interface{}

	// Pretty says to make the HTML look pretty before outputting.
	Pretty bool

	// Componenet is used to indicate that this is a snippet of code, not a full document.
	// As such, <html> and <head> tags will be suppressed.
	Component bool

	pool sync.Pool

	tmpl *template.Template
}

// Compile compiles the Doc into a template for execution internally.  Must be called before Execute().
// You only need to do this once.
func (d *Doc) Compile() error {
	if err := d.validate(); err != nil {
		return err
	}

	if d.Head != nil {
		if err := d.Head.compile(); err != nil {
			return err
		}
	}

	if d.Component {
		if d.Body != nil {
			d.Body.Component = true
		}
	}
	if err := d.Body.compile(); err != nil {
		return err
	}

	doc, err := template.New("doc").Parse(docTmpl)
	if err != nil {
		return err
	}
	d.tmpl = doc
	d.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	return nil
}

// Execute executes the internal template and writes the output to the io.Writer.
func (d *Doc) Execute(w io.Writer, data interface{}) error {
	if d.tmpl == nil {
		return fmt.Errorf("must call Compile() before Execute() on Doc type")
	}

	if d.Pretty {
		return d.tmpl.ExecuteTemplate(gohtml.NewWriter(w), "doc", pipeline{Self: d, Data: data})
	}

	return d.tmpl.ExecuteTemplate(w, "doc", pipeline{Self: d, Data: data})
}

// Render calls execute execute and returns the string value.
func (d *Doc) Render(data interface{}) (template.HTML, error) {
	w := d.pool.Get().(*strings.Builder)
	defer d.pool.Put(w)
	w.Reset()

	if err := d.Execute(w, data); err != nil {
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

// Element represents an HTML 5 element.
type Element interface {
	Execute(data interface{}) template.HTML
	compile() error
	isElement()
}

type outputAble interface {
	outputAble()
	fmt.Stringer
}

// DynamicFunc is a function that uses dynamic server data to return Elements that will be rendered.
type DynamicFunc func() []Element

type dynamic struct {
	f    DynamicFunc
	pool sync.Pool
}

func (d *dynamic) Execute(data interface{}) template.HTML {
	buff := d.pool.Get().(*strings.Builder)
	defer d.pool.Put(buff)
	buff.Reset()

	for _, e := range d.f() {
		buff.WriteString(string(e.Execute(data)))
	}
	return template.HTML(buff.String())
}

func (d *dynamic) isElement() {}

func (d *dynamic) compile() error {
	return nil
}

// Dynamic wraps a DynamicFunc so that it implements Element.
func Dynamic(f DynamicFunc) Element {
	return &dynamic{f: f}
}

// TextElement is an element that represents text, usually in a value. It is not valid everywhere.
type TextElement string

func (t TextElement) isElement() {}

func (t TextElement) Execute(data interface{}) template.HTML {
	return template.HTML(t)
}

func (t TextElement) compile() error {
	return nil
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
		if name == "TagValue" {
			continue
		}

		field := val.Field(i)

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

		if suffix := sf.Tag.Get("suffix"); suffix != "" {
			str = str + suffix
		}

		if isNaked {
			out = append(out, strings.ToLower(name))
		} else {
			out = append(out, fmt.Sprintf("%s=%q", strings.ToLower(name), str))
		}
	}

	return strings.Join(out, " ")
}
