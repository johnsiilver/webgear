// +build js wasm

/*
Package wasm provides the application entrance for WASM applications and the API for making UI changes.  
In this context, WASM applications are built off the webgear frameworks and update via this package's
UI object.

This package is used in two different places
	Inside your Wasm app where you will use the Wasm type to create your client.
	Inside your web server that is serving the Wasm app where you will use the Handler() func to generate a handler that loads your app.

WASM apps are binaries that are pre-compiled but loaded via JS served from a web browser. Normally it requires
an HTML page and JS bootstrapping. We provide the http.Handler that can mount this app and all the 
bootstrapping JS code.

This package along with the webgear/html and webgear/component remove the need to use HTML or JS within
your application.

The app handles all the client display while and should fetch content from the server in which to display.
This should be done via REST or other HTTP based calls.
*/
package wasm

import (
	"fmt"
	"bytes"
	"sync"
	"context"
	"html/template"
	"net/http"
	"net/url"
	"syscall/js"
	"reflect"

	"github.com/johnsiilver/webgear/html"
	
	"github.com/ulule/deepcopier"
)

type buffPool struct {
	pool sync.Pool
}

func (b *buffPool) get() *bytes.Buffer {
	buff := b.pool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff
}

func (b *buffPool) put(buff *bytes.Buffer) {
	b.pool.Put(buff)
}

var bufferPool = &buffPool{
	pool: sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	},
}

// Wasm represents a self contained WASM application.
type Wasm struct {
	doc *html.Doc
	renderIn chan *html.Doc 
	updateMu sync.Mutex
}

// New creates a new Wasm instance.  The doc passed will be used to replace the contents of the document once our app is
// downloaded from the server.
func New(doc *html.Doc) *Wasm {
	if doc == nil {
		panic("wasm.New() has nil doc argument")
	}
	if doc.Body == nil {
		panic("wasm.New() has nil doc.Body argument")
	}

	w := &Wasm{
		doc: doc, 
		renderIn: make(chan *html.Doc, 1),
	}

	return w
}

// Run executes replaces the content of the current document with the *html.Doc passed in New(). This never returns unless
// the initial doc cannot be rendered, which will cause a panic.
func (w *Wasm) Run(ctx context.Context) {
	req, _ := http.NewRequestWithContext(ctx, "POST", "/", bytes.NewBuffer([]byte{}))
	
	buff := bufferPool.get()
	defer bufferPool.put(buff)

	err := w.doc.Execute(ctx, buff, req)
	if err != nil {
		panic(err)
	}

	js.Global().Get("document").Call("innerHTML", buff.String())

	select{}
}

// UI creates a new UI object for changing the current UI output. This call is thread-safe, but blocks on
// all future calls until *UI.Closed() is called on the returned object.
func (w *Wasm) UI() *UI {
	w.updateMu.Lock() // *UI.Closed() unlocks this.  Yeah, yeah, I KNOW!!!!!

	return newUI(w.doc.Body, w)
}

// Func is a WASM function that is passed the Javascript "this" and any arguments. As this is attached to events, the return value
// is never used and is here only to satisfy a type in the syscall/JS library.
type Func func(this js.Value, args []js.Value) interface{}

// Attach attaches a Func to an html.Element for a specific event like a mouse click. If release is set, the function will
// release its memory after being called. This should be used when something like a button will not be used again in this
// evocation. Attach spins out a new goroutine for you to prevent blocking calls from pausing the event loop. Note that this
// may have negative consequences on performance, but that pause is such a hastle to track down for every new dev that
// I don't care about that. Finally, this attaches based on the element's .GlobalAttrs.ID. If that is not set, this is going to
// panic.
func Attach(event html.EventType, element html.Element, release bool, fn Func) html.Element{
	var cb js.Func
	js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			go func() {
				if release {
					defer cb.Release()
				}
				fn(this, args)
			}()
			return nil
		},
	)
	id := reflect.ValueOf(element).FieldByName("GlobalAttrs").Interface().(html.GlobalAttrs).ID
	if id == "" {
		panic(fmt.Errorf("cannot assign event(%s) to an element(%T) that has no ID set", event, element))
	}
	js.Global().Get("document").Call("getElementById", id).Call("addEventListener", string(event), fn)
	return element
}

// handler implements http.Handler by serving up an *html.Doc.
type handler struct {
	doc *html.Doc
}

// ServeHTTP implements http.Handler.ServeHTTP.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.doc.Execute(r.Context(), w, r)
}

// Handler will return a Handler that will return code to laod your
func Handler(u *url.URL) (http.Handler, error) {
	p := u.String()
	if p == "" {
		return nil, fmt.Errorf("the url passed(%s) to wasm.Handler was invalid", p)
	}

	doc := &html.Doc{
		Head: &html.Head{
			Elements: []html.Element{
				&html.Meta{Charset: "UTF-8"},
				&html.Script{TagValue: template.JS(wasmExec)},
				&html.Script{
					TagValue: template.JS(
						fmt.Sprintf(
`
const go = new Go();
WebAssembly.instantiateStreaming(fetch("%s"), go.importObject).then((result) => {
	go.run(result.instance);
});
`, p),
					),
				},
			},
		},
	}

	return handler{doc: doc}, nil
}

func copyBody(body *html.Body) *html.Body {
	cp := &html.Body{}
	deepcopier.Copy(body).To(cp)
	return cp
}