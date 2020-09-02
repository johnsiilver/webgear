// +build js,wasm

package html

import (
	"fmt"
	"reflect"
	"strings"
)

type IDRoot map[string]*Node

// element returns element at path. So if we are id "world" inside Component "hello", we would pass
// []string{"hello", "world"}.
func (i IDRoot) element(path []string) (Element, error) {
	root := i
	for i, p := range path[:len(path)-1] {
		v, ok := root[p]
		if !ok {
			root = v.Component
			continue
		}
		if v.Component == nil {
			return nil, fmt.Errorf("looking up element(%s) %q does not exist", strings.Join(path, "."), path[i])
		}
		root = v.Component
	}
	v, ok := root[path[len(path)-1]]
	if !ok {
		return nil, fmt.Errorf("looking up element(%s) %q does not exist", strings.Join(path, "."), path[len(path)-1])
	}
	if v.ElementNode.Element == nil {
		return nil, fmt.Errorf("looking up element(%s) %q does not exist", strings.Join(path, "."), path[len(path)-1])
	}
	return v.ElementNode.Element, nil
}

type Node struct {
	ElementNode *ElemNode
	Component   IDRoot
	ShadowPath  []string
}

func (n *Node) isComponent() bool {
	if n.Component == nil {
		return false
	}
	return true
}

func (n *Node) isElement() bool {
	if n.ElementNode.Element == nil {
		return false
	}
	return true
}

type ElemNode struct {
	Element Element
	Parent  Element
	// Path is the full path to this node: aka shadowPath + id.
	Path []string
}

func (e *ElemNode) ID() string {
	return GetElementID(e.Element)
}

func (e *ElemNode) ParentID() string {
	return GetElementID(e.Parent)
}

func replaceElementInNode(parent Element, id string, element Element) error {
	eID := GetElementID(element)
	if eID == "" {
		return fmt.Errorf("cannot replace an Element if the Element does not have an ID")
	}
	var pID string
	if _, ok := parent.(*Body); ok {
		pID = "body"
	} else {
		pID = GetElementID(parent)
	}
	if pID == "" {
		return fmt.Errorf("cannot replace an Element if the parent Element does not have an ID or isn't of type *Body")
	}

	if i, ok := element.(Initer); ok {
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
			sliceElem := slice.Index(i).Interface().(Element)
			if GetElementID(sliceElem) == id {
				slice.Index(i).Set(reflect.ValueOf(element))
				return nil
			}
		}
		return fmt.Errorf("couldn't find element(%s) in parent(%s)", eID, pID)
	}
	return nil
}

func addElementToNode(node Element, element Element) error {
	eID := GetElementID(element)
	if eID == "" {
		return fmt.Errorf("cannot add an Element to a node if the Element does not have an ID")
	}
	nID := GetElementID(node)
	if nID == "" {
		return fmt.Errorf("cannot add an Element if the parent Element does not have an ID")
	}

	if i, ok := element.(Initer); ok {
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

func deleteElement(id string, m IDRoot) (parent Element, err error) {
	n, ok := m[id]
	if !ok {
		return nil, fmt.Errorf("id %q could not be found", id)
	}

	parent = n.ElementNode.Parent
	if GetElementID(parent) == "" {
		return nil, fmt.Errorf("can't delete node(%s) whose parent does not have an ID", id)
	}

	pval := removePtr(reflect.ValueOf(parent))
	field := pval.FieldByName(elementField)
	if field.IsValid() {
		pval.FieldByName(elementField).Set(reflect.ValueOf(nil))
		return parent, nil
	}

	slice := pval.FieldByName(elementsField)
	if slice.IsValid() {
		newSlice := reflect.MakeSlice(reflect.TypeOf([]Element{}), 0, slice.Len()-1)
		for i := 0; i < slice.Len(); i++ {
			field := slice.Index(i)
			if GetElementID(field.Interface().(Element)) == id {
				continue
			}
			newSlice = reflect.Append(newSlice, field)
		}
		pval.FieldByName(elementsField).Set(newSlice)
		delete(m, id)
		return parent, nil
	}
	return nil, fmt.Errorf("can't delete element ID(%s): bug: parent(%T) doesn't seem to have .Element or .Elememnts attribute that contain it", id, parent)
}

const wasmUpdatedAttr = "XXXWasmUpdated"

// setUpdateFlag sets an elements update flag to true. Will error if the type does not have the flag.
func setUpdateFlag(e Element) error {
	if v, ok := e.(GearType); ok {
		v.SetUpdateFlag()
		return nil
	}
	val := removePtr(reflect.ValueOf(e))
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("element %T does not contain a GlobalAttrs attribute", e)
	}

	ga := val.FieldByName(globalAttrsField)
	if !ga.IsValid() {
		return fmt.Errorf("element %T does not contain a GlobalAttrs attribute", e)
	}

	ga.FieldByName(wasmUpdatedAttr).SetBool(true)
	val.FieldByName(globalAttrsField).Set(ga)
	return nil
}

func removeUpdateFlag(e Element) {
	if v, ok := e.(GearType); ok {
		v.RemoveUpdateFlag()
		return
	}
	val := removePtr(reflect.ValueOf(e))
	if val.Kind() != reflect.Struct {
		return
	}

	ga := val.FieldByName(globalAttrsField)
	if !ga.IsValid() {
		return
	}

	ga.FieldByName(wasmUpdatedAttr).SetBool(false)
	val.FieldByName(globalAttrsField).Set(ga)
	return
}

func getUpdate(e Element) bool {
	if v, ok := e.(GearType); ok {
		return v.UpdateFlag()
	}
	val := removePtr(reflect.ValueOf(e))
	if val.Kind() != reflect.Struct {
		return false
	}

	ga := val.FieldByName(globalAttrsField)
	if !ga.IsValid() {
		return false
	}

	return ga.FieldByName(wasmUpdatedAttr).Interface().(bool)
}

func removePtr(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}
