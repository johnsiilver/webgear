// +build js,wasm

package html

import (
	"fmt"
	"log"
	"strings"
	"syscall/js"
)

/*
// NewUI constructs a new UI for altering items within the scope of the attached element.
// Design Note: this is passed instead of *UI to do lazy initialization avoiding the setup
// cost if this would not be used.
type NewUI func() *UI


// UI creates a new UI object for changing the current UI output. This call is thread-safe, but blocks on
// all future calls until *UI.Closed() is called on the returned object.
func UI() *UI {
	w.updateMu.Lock() // *UI.Closed() unlocks this.  Yeah, yeah, I KNOW!!!!!

	return newUI(w.doc.Body, w)
}
*/

// WasmFunc is a Go function that is called by Javascript. WasmFunc is used to attach a Go function
// to events in WASM code. "this" is set to the js.Value of the object that the event was attached to
// in its state at the time of the call. root is set to the value of the root object that the "this"
// object is within. Without componenets, this will be "document". If contained in a set of components,
// this will be the component's shadowRoot that this element is contained within.
type WasmFunc func(this js.Value, root js.Value)

type wasmEvent struct {
	id            string
	handlerEvent  EventType
	listenerEvent ListenerType
	release       bool
	fn            WasmFunc
}

// Call calls the attached event passing along the component's shadowPath to allow finding the
// exact id we are looking for by querying shadowRoots.
func (w wasmEvent) Call(doc *Doc, shadowPath []string) {
	thisFunc := func() js.Value {
		v, err := elementByID(shadowPath, w.id)
		if err != nil {
			log.Printf("cannot get the value of element %s for a callback", w.id)
			return js.Undefined()
		}
		return v
	}
	rootFunc := func() js.Value {
		v, err := rootByPath(shadowPath)
		if err != nil {
			log.Printf("cannot get the value of a root element %s for a callback", strings.Join(shadowPath, "."))
			return js.Undefined()
		}
		return v
	}

	cb := w.makeCallback(thisFunc, rootFunc)

	log.Printf("ElementID(%s) ShadowPath(%v)", w.id, shadowPath)
	element, err := elementByID(shadowPath, w.id)
	if err != nil {
		log.Printf("cannot attach an event(id:%s)(shadowPath:%v): %s", w.id, shadowPath, err)
		return
	}
	if w.listenerEvent != "" {
		log.Printf("(%s).Call('addEventListener', '%s', <func>", w.id, w.listenerEvent)
	} else {
		log.Printf("(%s).Set('%s', <func>", w.id, w.handlerEvent)
	}

	log.Printf("the element we are attaching the event to: %s", element.Get("outerHTML"))
	if w.listenerEvent != "" {
		element.Call("addEventListener", string(w.listenerEvent), cb)
		log.Printf("added listener event to %q: type %s: %v", w.id, w.listenerEvent, cb)
		return
	}
	if w.handlerEvent != "" {
		element.Set(string(w.handlerEvent), cb)
		log.Printf("added handler event to %q: type %s: %v", w.id, w.handlerEvent, cb)
		return
	}
	panic("event attached for id(%s) with neither listenerEvent or handlerEvent")
}

/*
func (w wasmEvent) mapper(doc *Doc, shadowPath []string) (idRoot, error) {
	root := idRoot{}
	if err := mapIDs(doc.Body, root); err != nil {
		return nil, err
	}
	for _, p := range shadowPath {
		r, ok := root[p]
		if !ok {
			return nil, fmt.Errorf("could not find shadow path(%s) element %s", strings.Join(shadowPath, "."), p)
		}
		root = r.component
	}
	return root, nil
}

func (w wasmEvent) UI(mapper func() (idRoot, error)) *UI {
	root, err := mapper()
	if err != nil {
		log.Println("event on %q triggered UI() and mapper had error: %v", w.id, err)
		return nil
	}
}
*/

// makeCallback wraps the user defined WasmFunc inside the js.Func that is needed to embed the func
// in Javascript calls. "thisFunc" is called to get the current object the event is
// attached to. "rootFunc" returns the current root object the attached object is within. If not
// inside a component, the is "document". Otherwise it is the component.shadowRoot.
func (w wasmEvent) makeCallback(thisFunc, rootFunc func() js.Value) js.Func {
	var cb js.Func
	cb = js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			log.Println("event fired")
			go func() {
				if w.release {
					defer cb.Release()
				}
				w.fn(thisFunc(), rootFunc())
			}()
			return nil
		},
	)
	return cb
}

// AddWasmHandler causes a Go function to be bound to an event on this object. If release is true,
// the function is released (and unusable) after it is triggered. This is ignored if not
// using the wasm module. This is exposed to allow the Wasm module to access this. There
// is no compatibility promise for this method, you should use wasm.Attach().
func (e *Events) AddWasmHandler(id string, et EventType, fn WasmFunc, release bool) *Events {
	if id == "" {
		panic("AddWasmHandler cannot be called with id set to empty string")
	}
	if et == "" {
		panic("AddWasmHandler cannot get EventType empty string")
	}
	if fn == nil {
		panic("AddWasmHandler cannot receive a nil WasmFunc")
	}
	if e.wasmEvents == nil {
		e.wasmEvents = []wasmEvent{}
	}
	e.wasmEvents = append(
		e.wasmEvents.([]wasmEvent),
		wasmEvent{
			id:           id,
			handlerEvent: et,
			release:      release,
			fn:           fn,
		},
	)
	return e
}

// AddWasmListener causes a Go function to be bound to an event on this object. If release is true,
// the function is released (and unusable) after it is triggered. This is ignored if not
// using the wasm module. This is exposed to allow the Wasm module to access this. There
// is no compatibility promise for this method, you should use wasm.Attach().
func (e *Events) AddWasmListener(id string, et ListenerType, fn WasmFunc, release bool) *Events {
	if id == "" {
		panic("AddWasm cannot be called with id set to empty string")
	}
	if et == "" {
		panic("AddWasm cannot get EventType empty string")
	}
	if fn == nil {
		panic("AddWasm cannot receive a nil WasmFunc")
	}
	if e.wasmEvents == nil {
		e.wasmEvents = []wasmEvent{}
	}
	e.wasmEvents = append(
		e.wasmEvents.([]wasmEvent),
		wasmEvent{
			id:            id,
			listenerEvent: et,
			release:       release,
			fn:            fn,
		},
	)
	return e
}

// WasmEvents returns a list of functions that will attach Wasm events specified by
// .AddWasm(). This is for internal use only and has no compatibility promises.
func (e *Events) WasmEvents() []func(*Doc, []string) {
	if e.wasmEvents == nil {
		return nil
	}
	l := []func(*Doc, []string){}
	for _, event := range e.wasmEvents.([]wasmEvent) {
		l = append(l, event.Call)
	}
	return l
}

// elementByID roots through an element's top level components until it can do a getElementByID on
// the containing shadowRoot and returns the value.
func elementByID(shadowPath []string, id string) (js.Value, error) {
	root, err := rootByPath(shadowPath)
	if err != nil {
		return js.Undefined(), err
	}

	element := root.Call("getElementById", id)
	if element.IsUndefined() {
		fullPath := strings.Join(shadowPath, ".shadowRoot.") + id
		return js.Undefined(), fmt.Errorf("elementByID(%s): %s was undefined", fullPath, id)
	}
	return element, nil
}

func rootByPath(shadowPath []string) (js.Value, error) {
	root := js.Global().Get("document")
	if !root.Truthy() {
		return js.Undefined(), fmt.Errorf("the root document was undefined/null")
	}
	for _, component := range shadowPath {
		log.Printf("running getElementById(%s): %v", component, root.Call("getElementById", component))
		root = root.Call("getElementById", component).Get("shadowRoot")
		if !root.Truthy() {
			fullPath := strings.Join(shadowPath, ".shadowRoot.")
			return js.Undefined(), fmt.Errorf("rootByPath(%s): component was undefined", fullPath)
		}
	}
	return root, nil
}
