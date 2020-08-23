/*
Package builder provides a programatic way of building the HTML tree structure in a way that can be more easily read
without the documenet structure becoming a right leaning mess, such as:

	&Doc{
		Head: &Head{},
		Body: &Body{
			Elements: []Element{
				&Div{
					Elements: []Element{
						&Table{
							Elements: []TableElement{
								&TR{
									Elements: []TRElement{
										&TD{},
										&TD{},
									},
								},
							},
						},
					},
				},
			},
		},
	}

With our builder, this becomes
	b := NewHTML(&Head{}, &Body{})

	b.Into(&Div{}) // Adds the div and moves into the div
	b.Into(&Table{}) // Adds the table and moves into the table
	b.Into(&TR{}) // Adds the table row and moves into the table row
	b.Add(&TD{}, &TD{})  // Adds two table role elements, but stays in the row.
	b.Up() // We now move back to the table, if we called b.Up() again, we'd be at the div

*/
package builder

import (
	"fmt"
	"reflect"

	"github.com/johnsiilver/webgear/html"
)

type node struct {
	e html.Element
	p *node
}

func upNode(node *node) *node {
	if node == nil {
		panic("cannot go up a nil node")
	}

	return node.p
}

func addNode(add *node, to *node) *node {
	if add == nil {
		panic("cannot add a nil node")
	}
	if to == nil {
		return add
	}
	add.p = to
	return add
}

// HTML provides a builder for constucting HTML objects that can be rendered.  This attempts to allow tooling and other software
// constructs to be made that aren't right leaning pages of text.
type HTML struct {
	doc     *html.Doc
	current *node
}

// NewHTML creates a new HTML buider.
func NewHTML(head *html.Head, body *html.Body) *HTML {
	if head == nil || body == nil {
		panic("NewHTML called with a nil head or body arguement")
	}
	return &HTML{
		doc: &html.Doc{Head: head, Body: body},
	}
}

// Doc returns the *html.Doc.
func (h *HTML) Doc() *html.Doc {
	return h.doc
}

// AddHead adds elements to the doc's head.
func (h *HTML) AddHead(elements ...html.Element) {
	h.doc.Head.Elements = append(h.doc.Head.Elements, elements...)
}

// Add adds an element into the focus object, the object focus is not changed. If there is no object in focus
// this element is added to the doc's body.
func (h *HTML) Add(elements ...html.Element) {
	if h.current == nil {
		h.doc.Body.Elements = append(h.doc.Body.Elements, elements...)
		return
	}

	v := reflect.ValueOf(h.current.e).Elem().FieldByName("Elements")
	if !v.IsValid() {
		panic(fmt.Sprintf("(Add) cannot add to current element in builder.HTML(%T), it does not have an 'Elements' field", h.current.e))
	}

	for _, e := range elements {
		reflect.ValueOf(h.current.e).Elem().FieldByName("Elements").Set(reflect.Append(v, reflect.ValueOf(e)))
	}
}

// Into inserts element into the current focus object and then changes the focus object to this element. It returns the
// focus element.
func (h *HTML) Into(element html.Element) html.Element {
	parent := h.current
	h.current = &node{e: element, p: parent}

	if parent == nil {
		h.doc.Body.Elements = append(h.doc.Body.Elements, element)
		return element
	}

	v := reflect.ValueOf(parent.e).Elem().FieldByName("Elements")
	if !v.IsValid() {
		panic(fmt.Sprintf("(Into) cannot add to current element in builder.HTML(%T), it does not have an 'Elements' field", parent.e))
	}

	reflect.ValueOf(parent.e).Elem().FieldByName("Elements").Set(reflect.Append(v, reflect.ValueOf(element)))
	return h.current.e
}

// Up changes the current focus object to that object's parent. If there is no parent then the focus will become the
// doc's body. If the focus is already the doc's body, this will panic. Returns the HTML objects so you can call Up().Up() .
func (h *HTML) Up() *HTML {
	h.current = upNode(h.current)
	return h
}
