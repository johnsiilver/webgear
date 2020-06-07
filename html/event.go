package html

import (
	"fmt"
	"html/template"
	"log"
	"strings"
)

type event struct {
	key   string
	value string
}

func (e *event) String() string {
	if e.value == "" {
		log.Printf("an %q event was provided an empty scriptName, skipping", e.key)
		return ""
	}
	return fmt.Sprintf("%s=%q", e.key, e.value)
}

// Events represents an HTML event that triggers a javascript function.
// This is not used directly, but accessible via HTML elements.
// Once used in by an Execute(), the output will always be the same regardless of changes.
type Events struct {
	events []event

	str     string
	builder strings.Builder
}

func (e *Events) Attr() template.HTMLAttr {
	if e == nil {
		return ""
	}
	if e.str != "" {
		return template.HTMLAttr(e.str)
	}

	b := strings.Builder{}

	for i, event := range e.events {
		if i+1 == len(e.events) {
			b.WriteString(event.String())
		} else {
			b.WriteString(event.String() + " ")
		}
	}
	return template.HTMLAttr(b.String())
}

func (e *Events) OnAfterPrint(scriptName string) *Events {
	e.events = append(e.events, event{"onafterprint", scriptName})
	return e
}

func (e *Events) OnBeforePrint(scriptName string) *Events {
	e.events = append(e.events, event{"onbeforeprint", scriptName})
	return e
}

func (e *Events) OnBeforeUnload(scriptName string) *Events {
	e.events = append(e.events, event{"onbeforeunload", scriptName})
	return e
}

func (e *Events) OnError(scriptName string) *Events {
	e.events = append(e.events, event{"onerror", scriptName})
	return e
}

func (e *Events) OnHashChange(scriptName string) *Events {
	e.events = append(e.events, event{"onhashchange", scriptName})
	return e
}

func (e *Events) OnLoad(scriptName string) *Events {
	e.events = append(e.events, event{"onload", scriptName})

	return e
}
func (e *Events) OnMessage(scriptName string) *Events {
	e.events = append(e.events, event{"onmessage", scriptName})
	return e
}

func (e *Events) OnOffline(scriptName string) *Events {
	e.events = append(e.events, event{"onoffline", scriptName})
	return e
}

func (e *Events) OnOnline(scriptName string) *Events {
	e.events = append(e.events, event{"ononline", scriptName})
	return e
}

func (e *Events) OnPageHide(scriptName string) *Events {
	e.events = append(e.events, event{"onpagehide", scriptName})
	return e
}

func (e *Events) OnPageShow(scriptName string) *Events {
	e.events = append(e.events, event{"onpageshow", scriptName})
	return e
}

func (e *Events) OnPopState(scriptName string) *Events {
	e.events = append(e.events, event{"onpopstate", scriptName})
	return e
}

func (e *Events) OnResize(scriptName string) *Events {
	e.events = append(e.events, event{"onresize", scriptName})
	return e
}

func (e *Events) OnStorage(scriptName string) *Events {
	e.events = append(e.events, event{"onstorage", scriptName})
	return e
}

func (e *Events) OnUnload(scriptName string) *Events {
	e.events = append(e.events, event{"onunload", scriptName})
	return e
}

// OnBlur is the script that fires the moment that the element loses focus.
func (e *Events) OnBlur(scriptName string) *Events {
	e.events = append(e.events, event{"onblur", scriptName})
	return e
}

// OnChange is the script that fires the moment when the value of the element is changed.
func (e *Events) OnChange(scriptName string) *Events {
	e.events = append(e.events, event{"onchange", scriptName})
	return e
}

// OnContextMenu is the script to be run when a context menu is triggered.
func (e *Events) OnContextMenu(scriptName string) *Events {
	e.events = append(e.events, event{"oncontextmenu", scriptName})
	return e
}

// OnFocus is the script that fires the moment when the element gets focus.
func (e *Events) OnFocus(scriptName string) *Events {
	e.events = append(e.events, event{"onfocus", scriptName})
	return e
}

// OnInput is the script to be run when an element gets user input.
func (e *Events) OnInput(scriptName string) *Events {
	e.events = append(e.events, event{"oninput", scriptName})
	return e
}

// OnInvalid is the script to be run when an element is invalid.
func (e *Events) OnInvalid(scriptName string) *Events {
	e.events = append(e.events, event{"oninvalid", scriptName})
	return e
}

// OnReset is the script that fires when the Reset button in a form is clicked.
func (e *Events) OnReset(scriptName string) *Events {
	e.events = append(e.events, event{"onreset", scriptName})
	return e
}

// OnSearch is the script that fires when the user writes something in a search field (for <input="search">).
func (e *Events) OnSearch(scriptName string) *Events {
	e.events = append(e.events, event{"onsearch", scriptName})
	return e
}

// OnSelect is the script that fires after some text has been selected in an element.
func (e *Events) OnSelect(scriptName string) *Events {
	e.events = append(e.events, event{"onselect", scriptName})
	return e
}

// OnSubmit is the script that fires when a form is submitted.
func (e *Events) OnSubmit(scriptName string) *Events {
	e.events = append(e.events, event{"onsubmit", scriptName})
	return e
}

// OnKeyDown is the script that fires when a user is presses a key.
func (e *Events) OnKeyDown(scriptName string) *Events {
	e.events = append(e.events, event{"onkeydown", scriptName})
	return e
}

// OnKeyPress is the script that fires when a user is pressing a key.
func (e *Events) OnKeyPress(scriptName string) *Events {
	e.events = append(e.events, event{"onkeypress", scriptName})
	return e
}

// OnKeyUp is the script that fires when a user releases a key.
func (e *Events) OnKeyUp(scriptName string) *Events {
	e.events = append(e.events, event{"onkeyup", scriptName})
	return e
}

// OnClick is the script that fires on a mouse click on the element.
func (e *Events) OnClick(scriptName string) *Events {
	e.events = append(e.events, event{"onclick", scriptName})
	return e
}

// OnDblClick is the script that fires on a mouse double-click on the element.
func (e *Events) OnDblClick(scriptName string) *Events {
	e.events = append(e.events, event{"ondblclick", scriptName})
	return e
}

// OnMouseMove is the script that fires when a mouse button is pressed down on an element.
func (e *Events) OnMouseMove(scriptName string) *Events {
	e.events = append(e.events, event{"onmousedown", scriptName})
	return e
}

// OnMouseOut is the script that fires when the mouse pointer moves out of an element.
func (e *Events) OnMouseOut(scriptName string) *Events {
	e.events = append(e.events, event{"onmouseout", scriptName})
	return e
}

// OnMouseOver is the script that fires when the mouse pointer moves over an element.
func (e *Events) OnMouseOver(scriptName string) *Events {
	e.events = append(e.events, event{"onmouseover", scriptName})
	return e
}

// OnMouseUp is the script that fires when a mouse button is released over an element.
func (e *Events) OnMouseUp(scriptName string) *Events {
	e.events = append(e.events, event{"onmouseup", scriptName})
	return e
}

// OnWheel is the script that fires when the mouse wheel rolls up or down over an element.
func (e *Events) OnWheel(scriptName string) *Events {
	e.events = append(e.events, event{"onwheel", scriptName})
	return e
}

/*
OnDrag 	script 	Script to be run when an element is dragged
OnDragEnd 	script 	Script to be run at the end of a drag operation
ondragenter 	script 	Script to be run when an element has been dragged to a valid drop target
ondragleave 	script 	Script to be run when an element leaves a valid drop target
ondragover 	script 	Script to be run when an element is being dragged over a valid drop target
ondragstart 	script 	Script to be run at the start of a drag operation
ondrop 	script 	Script to be run when dragged element is being dropped
onscroll 	script 	Script to be run when an element's scrollbar is being scrolled
oncopy 	script 	Fires when the user copies the content of an element
oncut 	script 	Fires when the user cuts the content of an element
onpaste 	script 	Fires when the user pastes some content in an element
onabort 	script 	Script to be run on abort
oncanplay 	script 	Script to be run when a file is ready to start playing (when it has buffered enough to begin)
oncanplaythrough 	script 	Script to be run when a file can be played all the way to the end without pausing for buffering
oncuechange 	script 	Script to be run when the cue changes in a <track> element
ondurationchange 	script 	Script to be run when the length of the media changes
onemptied 	script 	Script to be run when something bad happens and the file is suddenly unavailable (like unexpectedly disconnects)
onended 	script 	Script to be run when the media has reach the end (a useful event for messages like "thanks for listening")
onerror 	script 	Script to be run when an error occurs when the file is being loaded
onloadeddata 	script 	Script to be run when media data is loaded
onloadedmetadata 	script 	Script to be run when meta data (like dimensions and duration) are loaded
onloadstart 	script 	Script to be run just as the file begins to load before anything is actually loaded
onpause 	script 	Script to be run when the media is paused either by the user or programmatically
onplay 	script 	Script to be run when the media is ready to start playing
onplaying 	script 	Script to be run when the media actually has started playing
onprogress 	script 	Script to be run when the browser is in the process of getting the media data
onratechange 	script 	Script to be run each time the playback rate changes (like when a user switches to a slow motion or fast forward mode)
onseeked 	script 	Script to be run when the seeking attribute is set to false indicating that seeking has ended
onseeking 	script 	Script to be run when the seeking attribute is set to true indicating that seeking is active
onstalled 	script 	Script to be run when the browser is unable to fetch the media data for whatever reason
onsuspend 	script 	Script to be run when fetching the media data is stopped before it is completely loaded for whatever reason
ontimeupdate 	script 	Script to be run when the playing position has changed (like when the user fast forwards to a different point in the media)
onvolumechange 	script 	Script to be run each time the volume is changed which (includes setting the volume to "mute")
onwaiting 	script 	Script to be run when the media has paused but is expected to resume (like when the media pauses to buffer more data)
*/
