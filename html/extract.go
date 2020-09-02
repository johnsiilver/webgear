package html

import (
	"context"
	"log"
	"reflect"
)

// ExtractEvents extracts an *Events object from an Element. If there is no Events, this will be nil.
func ExtractEvents(element Element) *Events {
	val := reflect.ValueOf(element)

	// If it is *struct, get the struct and assign back to val.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	log.Println("ElementID: ", GetElementID(element))

	ga := val.FieldByName("Events")
	if !ga.IsValid() || ga.IsNil() {
		return nil
	}
	log.Println("Events is: ", ga.Kind())
	return ga.Interface().(*Events)
}

// GetElementID will return the Element's GlobalAttr.ID if it has one. Empty string if not.
// If the Element is a *Gear, Gear.GearID().
func GetElementID(e Element) string {
	if g, ok := e.(GearType); ok {
		return g.GearID()
	}
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

const (
	globalAttrsField = "GlobalAttrs"
	idField          = "ID"
	elementField     = "Element"
	elementsField    = "Elements"
	docField         = "Doc"
	gearsField       = "Gears"
)

// Walked is returned by Walker which walks all Elements from some root.
type Walked struct {
	// Element is the HTML Element.
	Element Element
	// Parent is the HTML Element parent of Element.
	Parent Element
	// ShadowPath is the path of component shadowRoots between the Walker "root" and this Element, not including it.
	ShadowPath []string
}

// Walker walks all elements from the root and returns them on the returned channel, including the root element.
// The passed context can be used to cancel the walker.
func Walker(ctx context.Context, root Element) chan Walked {
	ch := make(chan Walked, 1)
	go func() {
		defer close(ch)
		walk(ctx, root, nil, nil, ch)
	}()
	return ch
}

func walk(ctx context.Context, element Element, parent Element, shadowPath []string, ch chan Walked) {
	select {
	case <-ctx.Done():
		return
	case ch <- Walked{Element: element, Parent: parent, ShadowPath: shadowPath}:
	}

	// If the element we have is Gear, then it begins the start of a new shadowRoot.
	// We require all Gear's to have names, so we can start tracking the shadowRoots
	// to allow attached events in Wasm to work.
	if v, ok := element.(GearType); ok {
		n := make([]string, len(shadowPath)+1)
		copy(n, shadowPath)
		n[len(shadowPath)] = string(v.TagType())
		shadowPath = n
	}

	for _, child := range childElements(element, true) {
		walk(ctx, child, element, shadowPath, ch)
	}
}

func childElements(element Element, gearDescend bool) []Element {
	val := reflect.ValueOf(element)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	// Some types have a single child called "Element".
	field := val.FieldByName(elementField)
	if field.IsValid() {
		return []Element{field.Interface().(Element)}
	}

	// Most types have children called "Elements".
	field = val.FieldByName(elementsField)
	if field.IsValid() {
		if v, ok := field.Interface().([]Element); ok {
			return v
		}

		// This happens when we have things like []TableElement instead of []Element, so we need to repack it.
		elements := make([]Element, 0, field.Len())
		for i := 0; i < field.Len(); i++ {
			elements = append(elements, field.Index(i).Interface().(Element))
		}
		return elements
	}

	// If this thing is a Gear, grab the children.
	if gearDescend {
		children := []Element{}
		descendGear(val, &children)
		if len(children) > 0 {
			return children
		}
	}
	return nil
}

func descendGear(possibleGear reflect.Value, children *[]Element) {
	if possibleGear.Kind() == reflect.Ptr {
		possibleGear = possibleGear.Elem()
	}

	if !hasFields(possibleGear, docField, gearsField) {
		return
	}

	gear := possibleGear // We know its a gear now.

	doc := gear.FieldByName(docField)
	if !doc.IsNil() {
		body := doc.Elem().FieldByName("Body")
		if !body.IsNil() {
			*children = append(*children, childElements(body.Interface().(Element), true)...)
		}
	}

	gears := gear.FieldByName(gearsField)
	if !gears.IsNil() && gears.Len() > 0 {
		for i := 0; i < gears.Len(); i++ {
			*children = append(*children, gears.Index(i).Interface().(Element))
		}
	}
}

func hasFields(val reflect.Value, fieldNames ...string) bool {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return false
	}
	for _, fieldName := range fieldNames {
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			return false
		}
	}
	return true
}
