package wasm

import (
	"fmt"
	"reflect"

	"github.com/johnsiilver/webgear/html"
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

	ga := val.FieldByName("GlobalAttr")
	if !ga.IsValid() {
		return ""
	}
	return ga.FieldByName("ID").String()
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
	val := reflect.ValueOf(element)

	// If it is *struct, get the struct and assign back to val.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	ga := val.FieldByName("GlobalAttr")
	if !ga.IsValid() { // If it doesn't have GlobalAttr, it can't have children.
		return nil
	}

	id := ga.FieldByName("ID").String()
	if id != "" {
		if _, ok := m[id]; ok {
			return fmt.Errorf("multiple elements with ID %q, %T found and %T found", m[id], element)
		}
		m[id] = elemNode{element: element, parent: parent}
	}

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := val.Field(i)

		if t.Field(i).Anonymous || !field.CanInterface() {
			continue
		}

		switch real := field.Interface().(type) {
		case html.Element:
			if err := walkElement(element, real, m); err != nil {
				return err
			}
		case []html.Element:
			if err := mapIDs(element, real, m); err != nil {
				return err
			}
		}
	}
	return nil
}

func replaceElementInNode(parent html.Element, element html.Element) error {
	if i, ok := element.(html.Initer); ok {
		i.Init()
	}

	pval := reflect.ValueOf(parent)
	field := pval.FieldByName("Element")
	if field.IsValid() {
		pval.FieldByName("Element").Set(reflect.ValueOf(element))
		return nil
	}
	
	slice := pval.FieldByName("Elements")
	if slice.IsValid() {
		for i := 0; i < slice.Len(); i++ {
			sliceElem := slice.Index(i).Interface().(html.Element)
			if getElementID(sliceElem) == getElementID(element) {
				slice.Index(i).Set(reflect.ValueOf(element))
				return nil
			}
		}
		return fmt.Errorf("couldn't find element(%s) in parent(%s)", getElementID(element), getElementID(parent))
	}
	return nil
}

func addElementToNode(node html.Element, element html.Element) error {
	if i, ok := element.(html.Initer); ok {
		i.Init()
	}

	pval := reflect.ValueOf(node)
	field := pval.FieldByName("Element")
	if field.IsValid() {
		pval.FieldByName("Element").Set(reflect.ValueOf(element))
		return nil
	}
	
	slice := pval.FieldByName("Elements")
	if slice.IsValid() {
		reflect.Append(slice, reflect.ValueOf(element))
		return nil
	}
	return fmt.Errorf("can't add element(%T) to element(%T): node doesn't have .Element or .Elememnts attribute", element, node)
}

func deleteNode(id string, m map[string]elemNode) (parent html.Element, err error) {
	n, ok := m[id]
	if !ok {
		return nil, fmt.Errorf("id %q could not be found", id)
	}

	if getElementID(n.parent) == "" {
		return nil, fmt.Errorf("can't delete node(%s) whose parent does not have an ID", id)
	}

	pval := reflect.ValueOf(n.parent)
	field := pval.FieldByName("Element")
	if field.IsValid() {
		pval.FieldByName("Element").Set(reflect.ValueOf(nil))
		return n.parent, nil
	}

	slice := pval.FieldByName("Elements")
	if slice.IsValid() {
		newSlice := reflect.MakeSlice(reflect.TypeOf([]html.Element{}), 0, slice.Len()-1)
		for i := 0; i < slice.Len(); i++ {
			field := slice.Index(i)
			if getElementID(field.Interface().(html.Element)) == id {
				continue
			}
			reflect.Append(newSlice, field)
		}
		pval.FieldByName("Elements").Set(newSlice)
		return n.parent, nil
	}
	return nil, fmt.Errorf("can't delete element ID(%s): bug: parent(%T) doesn't seem to have .Element or .Elememnts attribute that contain it", id, n.parent)
}