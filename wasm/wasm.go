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
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sync"
	"syscall/js"
	"time"

	"github.com/johnsiilver/webgear/html"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

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
	doc      *html.Doc
	renderIn chan *html.Doc
	updateMu sync.Mutex

	ready   chan struct{}
	runMu   sync.Mutex
	running bool
}

// New creates a new Wasm instance.
func New() *Wasm {
	w := &Wasm{
		renderIn: make(chan *html.Doc, 1),
		ready:    make(chan struct{}),
	}

	return w
}

// SetDoc sets the initial doc that will be displayed.
func (w *Wasm) SetDoc(doc *html.Doc) {
	w.runMu.Lock()
	defer w.runMu.Unlock()
	if w.running {
		panic("cannot call Wasm.StartingDoc() once Run() has been called")
	}
	w.doc = doc
}

// Run executes replaces the content of the current document with the *html.Doc passed in New(). This never returns unless
// the initial doc cannot be rendered, which will cause a panic.
func (w *Wasm) Run(ctx context.Context) {
	w.runMu.Lock()
	if w.running {
		panic("Wasm.Run() already called")
	}
	w.running = true
	w.runMu.Unlock()

	if w.doc == nil {
		w.doc = &html.Doc{
			Head: &html.Head{},
			Body: &html.Body{},
		}
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", "/", bytes.NewBuffer([]byte{}))

	buff := bufferPool.get()
	defer bufferPool.put(buff)

	if w.doc.Head == nil {
		w.doc.Head = &html.Head{}
	}

	w.doc.Init()
	err := w.doc.Execute(ctx, buff, req)
	if err != nil {
		panic(err)
	}

	updater, err := html.NewDocUpdater(w.doc, docUpdaterHolder)
	if err != nil {
		panic(err)
	}

	//log.Println("document as it is in Run() before it is set: ", js.Global().Get("document").Get("documentElement").Get("outerHTML"))
	js.Global().Get("document").Call("open")
	js.Global().Get("document").Call("write", buff.String())
	// Setup event on the Body so that all other events will load when the body has loaded.
	w.initEvents()
	js.Global().Get("document").Call("close")

	log.Println("document as it set in Run():\n", js.Global().Get("document").Get("documentElement").Get("innerHTML"))
	docUpdaterHolder <- updater

	close(w.ready)
	select {}
}

// Ready will block until Run() has rendered the initial document. Mostly used in tests.
func (w *Wasm) Ready() {
	<-w.ready
}

// initEvents adds an event onto the window that loads all our other events
// once the body has finished loading.
func (w *Wasm) initEvents() {
	bodyID := w.doc.Body.ID
	bodyEvents := w.doc.Body.Events

	if bodyEvents == nil {
		w.doc.Body.Events = &html.Events{}
		bodyEvents = w.doc.Body.Events
	}
	if bodyID == "" {
		ga := w.doc.Body.GlobalAttrs
		ga.ID = "bodyID"
		bodyID = ga.ID
		w.doc.Body.GlobalAttrs = ga
	}

	fn := func(this js.Value, root js.Value, args interface{}) {
		w.doc.ExecuteDomCalls()
	}

	var cb js.Func
	cb = js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			log.Println("HEEEEEEEEEERRRRRRRRRRRRRRREEEEEEEEEEEEEEEEEE")
			go func() {
				wg := &sync.WaitGroup{}
				defer func () {
					wg.Wait()
					cb.Release()
				}()
				wg.Add(1)
				go func() {
					defer wg.Done()
					fn(js.Value{}, js.Value{}, nil)
				}()
			}()
			return nil
		},
	)

	//js.Global().Call("addEventListener", "load", cb)
	js.Global().Get("document").Call("addEventListener", "DOMContentLoaded", cb)
	//bodyEvents.AddWasmHandler(bodyID, html.OnLoad, fn, nil, true)
}

// AttachListener attaches a Func to an html.Element for a specific event like a mouse click (PLEASE READ THE REST BEFORE USE). If release is set, the function will
// release its memory after being called. This should be used when something like a button will not be used again in this
// evocation. Attach spins out a new goroutine for you to prevent blocking calls from pausing the event loop. Note that this
// may have negative consequences on performance, but that pause is such a hastle to track down for every new dev that
// I don't care about that. This func() attaches based on the element's .GlobalAttrs.ID. If that is not set, this is going to
// panic.
func AttachListener(event html.ListenerType, release bool, fn html.WasmFunc, args interface{}, element html.Element) html.Element {
	val := reflect.ValueOf(element)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic(fmt.Sprintf("Element of type %T does not have an ID field, which must be present to attach an event", element))
	}

	if !val.FieldByName("GlobalAttrs").IsValid() {
		panic(fmt.Errorf("cannot assign event(%s) to an element(%T) that has no ID set", event, element))
	}

	id := val.FieldByName("GlobalAttrs").Interface().(html.GlobalAttrs).ID
	if id == "" {
		panic(fmt.Errorf("cannot assign event(%s) to an element(%T) that has no ID set", event, element))
	}

	// If we don't have an Events on the Element, create it.
	eventsField := val.FieldByName("Events")
	if !eventsField.IsValid() {
		panic(fmt.Sprintf("Element type %T does not have an Events composition, so it cannot have Attach() called on it", element))
	}
	if eventsField.IsNil() {
		e := &html.Events{}
		eventsField.Set(reflect.ValueOf(e))
	}

	events := eventsField.Interface().(*html.Events)
	events.AddWasmListener(id, event, fn, args, release)
	return element
}

// AttachHandler attaches a Func to an html.Element for a specific event like a mouse click (PLEASE READ THE REST BEFORE USE). If release is set, the function will
// release its memory after being called. This should be used when something like a button will not be used again in this
// evocation. Attach spins out a new goroutine for you to prevent blocking calls from pausing the event loop. Note that this
// may have negative consequences on performance, but that pause is such a hastle to track down for every new dev that
// I don't care about that. This func() attaches based on the element's .GlobalAttrs.ID. If that is not set, this is going to
// panic. As this is an event handler, it replaces any existing handlers. Use even listeners to assign multiple actions.
func AttachHandler(event html.EventType, release bool, fn html.WasmFunc, args interface{}, element html.Element) html.Element {
	val := reflect.ValueOf(element)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic(fmt.Sprintf("Element of type %T does not have an ID field, which must be present to attach an event", element))
	}

	if !val.FieldByName("GlobalAttrs").IsValid() {
		panic(fmt.Errorf("cannot assign event(%s) to an element(%T) that has no ID set", event, element))
	}

	id := val.FieldByName("GlobalAttrs").Interface().(html.GlobalAttrs).ID
	if id == "" {
		panic(fmt.Errorf("cannot assign event(%s) to an element(%T) that has no ID set", event, element))
	}

	// If we don't have an Events on the Element, create it.
	eventsField := val.FieldByName("Events")
	if !eventsField.IsValid() {
		panic(fmt.Sprintf("Element type %T does not have an Events composition, so it cannot have Attach() called on it", element))
	}
	if eventsField.IsNil() {
		e := &html.Events{}
		eventsField.Set(reflect.ValueOf(e))
	}

	events := eventsField.Interface().(*html.Events)
	events.AddWasmHandler(id, event, fn, args, release)
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

// docUpdaterHolder holds the single DocUpdater.
var docUpdaterHolder = make(chan *html.DocUpdater, 1)

// GetDocUpdater gets a *DocUpdater for updating the DOM. There is only one of these available for
// use at any time. Use DocUpdater.Render() to release it. This will block if another call has
// not released it.
func GetDocUpdater() *html.DocUpdater {
	var du *html.DocUpdater
	ticker := time.NewTicker(10 * time.Second)
	for {
		ticker.Reset(10 * time.Second)
		select {
		case du = <-docUpdaterHolder:
		case <-ticker.C:
			log.Println("GetDocUpdater() called and hasn't returned after 10 seconds. " +
				"The last GetDocUpdater() caller didn't release it with a call to Render()")
			continue
		}
		break
	}
	ticker.Stop()

	return du
}
