// +build js,wasm

package html

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"syscall/js"
)

/*
This file holds Execute() and ExecuteAsGear() for all compilation targets that are wasm.
*/

// Execute executes the internal templates and writes the output to the io.Writer. This is thread-safe.
func (d *Doc) Execute(ctx context.Context, w io.Writer, r *http.Request) (err error) {
	if !d.initDone {
		return fmt.Errorf("Doc object did not have .Init() called before Execute()")
	}

	pipe := NewPipeline(ctx, r, w)
	pipe.Self = d

	if err := docTmpl.Execute(w, pipe); err != nil {
		return err
	}

	return pipe.HadError()
}

// ExecuteAsGear uses the Pipeline provided instead of creating one internally. This is for internal use only
// and no guarantees are made on its operation or that it will exist in the future. This is thread-safe.
func (d *Doc) ExecuteAsGear(pipe Pipeline) string {
	if !d.initDone {
		pipe.Error(fmt.Errorf("Doc object did not have .Init() called before Execute()"))
		return EmptyString
	}
	pipe.Self = d

	if err := docTmpl.Execute(pipe.W, pipe); err != nil {
		pipe.Error(err)
	}

	return EmptyString
}

// ExecuteDomCalls is run after the Doc has had Execute() or ExecuteAsGear() called in order
// to attach any of the events for WASM into the Dom.
func (d *Doc) ExecuteDomCalls() {
	log.Println("ExecuteDomCalls")

	for walked := range Walker(context.Background(), d.Body) {
		log.Printf("element: %T(%+v)", walked.Element, walked.ShadowPath)
		if events := ExtractEvents(walked.Element); events != nil {
			for _, fn := range events.WasmEvents() {
				fn(d, walked.ShadowPath)
			}
		}
	}
}

// DocUpdater is used do update the WASM's Doc representation and render the changes to the DOM.
// GetDocUpdater() is used to retrieve a DocUpater.
type DocUpdater struct {
	doc    *Doc
	holder chan *DocUpdater
	nodes  IDRoot
	mu     sync.Mutex
}

// NewDocUpdater is the constructor for DocUpdater.
func NewDocUpdater(doc *Doc, holder chan *DocUpdater) (*DocUpdater, error) {
	root := IDRoot{}
	if err := mapIDs(doc.Body, root); err != nil {
		return nil, err
	}
	return &DocUpdater{doc: doc, holder: holder, nodes: root}, nil
}

// NewFrom returns a new DocUpdater built from the internals of this one.
func (d *DocUpdater) NewFrom() (*DocUpdater, error) {
	return NewDocUpdater(d.doc, d.holder)
}

// UpdateElementByID updates a Doc Element by ID to now be the Element e.
func (d *DocUpdater) UpdateElementByID(id string, e Element) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	node, ok := d.nodes[id]
	if !ok {
		return fmt.Errorf("could not locate element ID(%s)", id)
	}
	if _, ok := node.ElementNode.Element.(*Component); ok {
		if _, ok := e.(GearType); ok {
			return fmt.Errorf("UpdateElementByID(%s) is trying to replace a Component tag with a Gear. " +
				"this is a mistake. You want to use id set to .GearID()")
		}
	}

	if err := setUpdateFlag(e); err != nil {
		return fmt.Errorf("UpdateElementByID(%s) had error: %w", id, err)
	}
	log.Printf("%T type", node.ElementNode.Parent)
	if err := replaceElementInNode(node.ElementNode.Parent, id, e); err != nil {
		return fmt.Errorf("UpdateElementByID(%s) had error: %w", id, err)
	}

	node.ElementNode.Element = e

	return nil
}

// AddChild updates an Element at toID to now have child Element e.
func (d *DocUpdater) AddChild(toID string, e Element) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	node, ok := d.nodes[toID]
	if !ok {
		return fmt.Errorf("could not locate element ID(%s)", toID)
	}
	if node.isComponent() {
		return fmt.Errorf("cannot call DocUpdater.AddChild(%s): is a Component", toID)
	}

	if err := setUpdateFlag(node.ElementNode.Element); err != nil {
		return fmt.Errorf("cannot call DocUpdater.AddChild(%s): %w", err)
	}

	if err := addElementToNode(node.ElementNode.Element, e); err != nil {
		return err
	}

	if id := GetElementID(e); id != "" {
		// (TODO): This really need to add more information and especially if e is a Gear.
		d.nodes[id] = &Node{ElementNode: &ElemNode{Element: e}}
	}
	return nil
}

// DeleteElementByID removes Element with ID id from the Doc.
func (d *DocUpdater) DeleteElementByID(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	parent, err := deleteElement(id, d.nodes)
	if err != nil {
		return err
	}

	if err := setUpdateFlag(parent); err != nil {
		return fmt.Errorf("cannot call DocUpdater.AddChild(%s): %w", err)
	}

	delete(d.nodes, id)

	return nil
}

// UpdateDOM renders the changes and releases the DocUpdater. This DocUpdater will panic
// if you attempt to use it after Render() is called.
func (d *DocUpdater) UpdateDOM() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	n, err := d.NewFrom()
	if err != nil {
		return err
	}

	defer func() {
		d.doc = nil
		d.nodes = nil
		d.holder <- n
		d.holder = nil
	}()

	buff := &bytes.Buffer{}
	for updateElement := range walkUpdatesOnly(d.doc) {
		log.Printf("ELEMENT THAT UPDATES: %s(%T)", GetElementID(updateElement.Element), updateElement.Element)
		if v, ok := updateElement.Element.(GearType); ok {
			if err := v.UpdateDOM(); err != nil {
				return err
			}
			continue
		}

		// Get the content.
		buff.Reset()
		pipe := NewPipeline(context.TODO(), &http.Request{}, buff)
		updateElement.Element.Execute(pipe)

		log.Println("element name to update: ", GetElementID(updateElement.Element))
		js.Global().Get("document").Call("getElementById", GetElementID(updateElement.Element)).Set("outerHTML", buff.String())
		removeUpdateFlag(updateElement.Element)
	}

	d.doc.ExecuteDomCalls()
	return nil
}

func mapIDs(root Element, m IDRoot) error {
	for walked := range Walker(context.Background(), root) {
		eid := GetElementID(walked.Element)
		if eid == "" {
			continue
		}
		// Its a root element.
		if len(walked.ShadowPath) == 0 {
			m[eid] = &Node{
				ElementNode: &ElemNode{
					Element: walked.Element,
					Parent:  walked.Parent,
					Path:    nil,
				},
			}
			continue
		}
		// Ok, find where to put it. If it has components above it that don't exist, add them.
		root := IDRoot{}
		for _, p := range walked.ShadowPath {
			v, ok := root[p]
			if ok {
				root = v.Component
				continue
			}
			// Doesn't do a composite literal because a node can both have an element and a Component
			// attribute, as a Gear is both.
			node := root[p]
			if node == nil {
				node = &Node{}
			}
			node.Component = IDRoot{}
			root[p] = node

			root = root[p].Component
		}
		// Now attach it.
		// Doesn't do a composite literal because a node can both have an element and a Component
		// attribute, as a Gear is both.
		node := root[eid]
		if node == nil {
			node = &Node{
				ElementNode: &ElemNode{},
			}
		}
		node.ElementNode.Element = walked.Element
		node.ElementNode.Parent = walked.Parent
		ns := make([]string, len(walked.ShadowPath), len(walked.ShadowPath)+1)
		copy(ns, walked.ShadowPath)
		node.ElementNode.Path = append(ns, eid)
		root[eid] = node
	}
	return nil
}

func walkUpdatesOnly(doc *Doc) chan Walked {
	ch := make(chan Walked, 1)
	go func() {
		defer close(ch)
		if getUpdate(doc.Body) {
			ch <- Walked{Element: doc.Body}
			return
		}
		walkUpdates(doc.Body, nil, ch)
	}()
	return ch
}

func walkUpdates(element Element, parent Element, ch chan Walked) {
	if getUpdate(element) {
		ch <- Walked{Element: element, Parent: parent}
		return
	}

	// If the element we have is Gear, then we do not descend.
	if _, ok := element.(GearType); ok {
		return
	}

	for _, child := range childElements(element, false) { // childElements in extract.go
		walkUpdates(child, element, ch)
	}
}
