package wasm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"syscall/js"

	"github.com/johnsiilver/webgear/html"

	"github.com/mohae/deepcopy"
)

// UI provides an object that controls updating the client's Web UI.
// This object allows you to stage a set of changes and then render that change on Close().
type UI struct {
	mu       sync.Mutex
	body     *html.Body
	nodes    map[string]elemNode
	validate []func() error
	updates  []func()
	wasm     *Wasm
	closed   bool
}

func newUI(body *html.Body, w *Wasm) *UI {
	nodes := map[string]elemNode{}
	if err := mapIDs(body, body.Elements, nodes); err != nil {
		panic(err)
	}

	return &UI{
		body:  copyBody(body),
		nodes: nodes,
		wasm:  w,
	}
}

// Close closes our our UI for updates and renders the changes. If already Closed this will panic.
// Once closed, any use of this object will cause a panic.
func (u *UI) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		panic("UI.Close() called on closed UI")
	}
	u.closed = true
	defer u.wasm.updateMu.Unlock() // Yuck

	for _, fn := range u.validate {
		if err := fn(); err != nil {
			return fmt.Errorf("an item that was supposed to update did not exist, no updates applied: %s", err)
		}
	}

	for _, fn := range u.updates {
		fn()
	}

	m := map[string]elemNode{}
	if err := mapIDs(u.body, u.body.Elements, m); err != nil {
		log.Fatalf("bad bug: %s", err)
	}
	u.wasm.doc.Body = u.body
	u.nodes = m
	return nil
}

// Body returns the current body object that is under transformation. This only represents the current
// live body object if no update calls have been made. This object is only safe to use until the next
// UI method call. The object should be considered read-only, manipulating the object can cause unknown
// issues.
func (u *UI) Body() *html.Body {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		panic("UI.Body() called on closed UI")
	}

	return u.body
}

// Update updates element with "id" to have the element "with".  "with" must have an ID.
func (u *UI) Update(id string, with html.Element) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		panic("UI.Update() called on closed UI")
	}

	if id == "" {
		return fmt.Errorf("Update() called with no ID set")
	}

	wid := getElementID(with)
	if wid == "" {
		return fmt.Errorf("Update() had ID %q with a replacement element with ID %q\n", id, getElementID(with))
	}

	if wid != id {
		if _, ok := u.nodes[wid]; ok {
			return fmt.Errorf("Update() was replacing element %q with elememnt with ID %q, however %q already exists", id, wid, wid)
		}
	}

	n, ok := u.nodes[id]
	if !ok {
		return fmt.Errorf("Update request received for ID %q, but that doesn't exist\n", id)
	}
	if fmt.Sprintf("%T", n.element) != fmt.Sprintf("%T", with) {
		return fmt.Errorf("Update() was updating a %T node, but used a %T node, which isn't allowed", n.element, with)
	}
	if err := replaceElementInNode(n.parent, id, with); err != nil {
		return fmt.Errorf("Update() was updating a %T node and got error: %w", n.parent, err)
	}
	buff := bufferPool.get()
	defer bufferPool.put(buff)

	pipe := html.NewPipeline(context.Background(), &http.Request{}, buff)
	with.Execute(pipe)
	u.validate = append(
		u.validate,
		func() error {
			log.Println("document text: ", js.Global().Get("document").Get("outerHTML"))
			log.Println("myDiv: ", js.Global().Get("document").Call("getElementById", "myDiv"))
			el := js.Global().Get("document").Call("getElementById", id)
			if !el.Truthy() {
				return fmt.Errorf("attempt to update element ID %q failed: element does not exist", id)
			}
			return nil
		},
	)
	u.updates = append(
		u.updates,
		func() {
			log.Printf("calling document.getElementById(%s).outerHTML = <stuff>", id)
			js.Global().Get("document").Call("getElementById", id).Set("outerHTML", buff.String())
		},
	)
	return nil
}

// AddTo adds the "with" element as a child of element with "id".
func (u *UI) AddTo(id string, with html.Element) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		panic("UI.AddTo() called on closed UI")
	}

	eid := getElementID(with)
	if eid == "" {
		return fmt.Errorf("AddTo(%s) passed element(%T) with no ID", id, with)
	}
	n, ok := u.nodes[id]
	if !ok {
		return fmt.Errorf("AddTo() call invalid: ID(%s) does not exist", id)
	}
	if err := addElementToNode(n.parent, with); err != nil {
		return fmt.Errorf("AddTo() request error: %w", err)
	}

	buff := bufferPool.get()
	defer bufferPool.put(buff)

	pipe := html.NewPipeline(context.Background(), &http.Request{}, buff)
	n.element.Execute(pipe)
	u.validate = append(
		u.validate,
		func() error {
			el := js.Global().Get("document").Call("getElementById", id)
			if !el.Truthy() {
				return fmt.Errorf("attempt to update element ID %q failed: element does not exist", id)
			}
			return nil
		},
	)
	u.updates = append(
		u.updates,
		func() {
			js.Global().Get("document").Call("getElementById", id).Set("innerHTML", buff.String())
		},
	)
	return nil
}

// Delete deletes an element with "id".
func (u *UI) Delete(id string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.closed {
		panic("UI.Delete() called on closed UI")
	}

	parent, err := deleteElement(id, u.nodes)
	if err != nil {
		return fmt.Errorf("Delete() error: %v", err)
	}

	buff := bufferPool.get()
	defer bufferPool.put(buff)

	pipe := html.NewPipeline(context.Background(), &http.Request{}, buff)
	parent.Execute(pipe)
	u.validate = append(
		u.validate,
		func() error {
			el := js.Global().Get("document").Call("getElementById", getElementID(parent))
			if !el.Truthy() {
				return fmt.Errorf("attempt to update element ID %q failed: element does not exist", id)
			}
			return nil
		},
	)
	u.updates = append(
		u.updates,
		func() {
			js.Global().Get("document").Call("getElementById", getElementID(parent)).Set("innerHTML", buff.String())
		},
	)
	return nil
}

func copyBody(body *html.Body) *html.Body {
	return deepcopy.Copy(body).(*html.Body)
}
