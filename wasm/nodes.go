package wasm

import (
	"fmt"
	"reflect"

	"github.com/johnsiilver/webgear/html"
)

const (
	globalAttrsField = "GlobalAttrs"
	idField = "ID"
	elementField = "Element"
	elementsField = "Elements"
)

type elemNode struct {
	element html.Element
	parent html.Element
}

func (e elemNode) ID() string {
	return getElementID(e.element)
}

func (e elemNode) ParentID() string {
	return getElementID(e.parent)
}

func getElementID(e html.Element) string {
	val := reflect.ValueOf(e)

	// If it is *struct, get the struct and assign back to val.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ""
	}

	ga := val.FieldByName(globalAttrsField)
	if !ga.IsValid() {
		return ""
	}
	return ga.FieldByName(idField).String()
}

// mapIDs maps all the IDs in an element
func mapIDs(parent html.Element, elements []html.Element, m map[string]elemNode) error {
	for _, element := range elements {
		if err := walkElement(parent, element, m); err != nil {
			return err
		}
	}
	return nil
}

func walkElement(parent, element html.Element, m map[string]elemNode) error {
	val := removePtr(reflect.ValueOf(element))

	if val.Kind() != reflect.Struct {
		return nil
	}

	ga := val.FieldByName(globalAttrsField)
	if !ga.IsValid() { // If it doesn't have GlobalAttr, it can't have children.
		return nil
	}

	id := ga.FieldByName(idField).String()
	if id != "" {
		if _, ok := m[id]; ok {
			return fmt.Errorf("multiple elements with ID %q, %T found and %T found", id, m[id], element)
		}
		m[id] = elemNode{element: element, parent: parent}
	}

	field := val.FieldByName(elementField)
	if field.IsValid() {
		if err := walkElement(element, field.Interface().(html.Element), m); err != nil {
			return err
		}
		return nil
	}
	field = val.FieldByName(elementsField)
	if err := mapIDs(element, field.Interface().([]html.Element), m); err != nil {
		return err
	}
	return nil
}

func replaceElementInNode(parent html.Element, id string, element html.Element) error {
	eID := getElementID(element)
	if eID == "" {
		return fmt.Errorf("cannot replace an Element if the Element does not have an ID")
	}
	pID := getElementID(parent)
	if pID == "" {
		return fmt.Errorf("cannot replace an Element if the parent Element does not have an ID")
	}

	if i, ok := element.(html.Initer); ok {
		i.Init()
	}

	pval := removePtr(reflect.ValueOf(parent))

	field := pval.FieldByName(elementField)
	if field.IsValid() {
		pval.FieldByName(elementField).Set(reflect.ValueOf(element))
		return nil
	}
	
	slice := pval.FieldByName(elementsField)
	if slice.IsValid() {
		for i := 0; i < slice.Len(); i++ {
			sliceElem := slice.Index(i).Interface().(html.Element)
			if getElementID(sliceElem) == id {
				slice.Index(i).Set(reflect.ValueOf(element))
				return nil
			}
		}
		return fmt.Errorf("couldn't find element(%s) in parent(%s)", eID, pID)
	}
	return nil
}

func addElementToNode(node html.Element, element html.Element) error {
	eID := getElementID(element)
	if eID == "" {
		return fmt.Errorf("cannot add an Element to a node if the Element does not have an ID")
	}
	nID := getElementID(node)
	if nID == "" {
		return fmt.Errorf("cannot add an Element if the parent Element does not have an ID")
	}

	if i, ok := element.(html.Initer); ok {
		i.Init()
	}

	pval := removePtr(reflect.ValueOf(node))
	if pval.Kind() != reflect.Struct {
		return fmt.Errorf("cannot add a child node to node type %T", node)
	}

	field := pval.FieldByName(elementField)
	if field.IsValid() {
		pval.FieldByName(elementField).Set(reflect.ValueOf(element))
		return nil
	}
	
	slice := pval.FieldByName(elementsField)
	if slice.IsValid() {
		slice = reflect.Append(slice, reflect.ValueOf(element))
		pval.FieldByName(elementsField).Set(slice)
		return nil
	}
	return fmt.Errorf("can't add element(%T) to element(%T): node doesn't have .Element or .Elements attribute", element, node)
}

func deleteElement(id string, m map[string]elemNode) (parent html.Element, err error) {
	n, ok := m[id]
	if !ok {
		return nil, fmt.Errorf("id %q could not be found", id)
	}

	if getElementID(n.parent) == "" {
		return nil, fmt.Errorf("can't delete node(%s) whose parent does not have an ID", id)
	}

	pval := removePtr(reflect.ValueOf(n.parent))
	field := pval.FieldByName(elementField)
	if field.IsValid() {
		pval.FieldByName(elementField).Set(reflect.ValueOf(nil))
		return n.parent, nil
	}

	slice := pval.FieldByName(elementsField)
	if slice.IsValid() {
		newSlice := reflect.MakeSlice(reflect.TypeOf([]html.Element{}), 0, slice.Len()-1)
		for i := 0; i < slice.Len(); i++ {
			field := slice.Index(i)
			if getElementID(field.Interface().(html.Element)) == id {
				continue
			}
			newSlice = reflect.Append(newSlice, field)
		}
		pval.FieldByName(elementsField).Set(newSlice)
		delete(m, id)
		return n.parent, nil
	}
	return nil, fmt.Errorf("can't delete element ID(%s): bug: parent(%T) doesn't seem to have .Element or .Elememnts attribute that contain it", id, n.parent)
}

func removePtr(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}