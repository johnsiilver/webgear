package html

import (
	"fmt"
	"html/template"
	"log"
	"strings"
)

// EventType represents a browser based event like a click, mouseover, etc...
type EventType string

const (
	OnAfterPrint   EventType = "onafterprint"
	OnBeforePrint  EventType = "onbeforeprint"
	OnBeforeUnload EventType = "onbeforeunload"
	OnError        EventType = "onerror"
	OnHashChange   EventType = "onhashchange"
	OnLoad         EventType = "onload"
	OnMessage      EventType = "onmessage"
	OnOffline      EventType = "onoffline"
	OnOnline       EventType = "ononline"
	OnPageHide     EventType = "onpagehide"
	OnPageShow     EventType = "onpageshow"
	OnPopState     EventType = "onpopstate"
	OnResize       EventType = "onresize"
	OnStorage      EventType = "onstorage"
	OnUnload       EventType = "onunload"
	OnBlur         EventType = "onblur"
	OnChange       EventType = "onchange"
	OnContextMenu  EventType = "oncontextmenu"
	OnFocus        EventType = "onfocus"
	OnInput        EventType = "oninput"
	OnInvalid      EventType = "oninvalid"
	OnReset        EventType = "onreset"
	OnSearch       EventType = "onsearch"
	OnSelect       EventType = "onselect"
	OnSubmit       EventType = "onsubmit"
	OnKeyDown      EventType = "onkeydown"
	OnKeyPress     EventType = "onkeypress"
	OnKeyUp        EventType = "onkeyup"
	OnClick        EventType = "onclick"
	OnDblClick     EventType = "ondblclick"
	OnMouseMove    EventType = "onmousemove"
	OnMouseOut     EventType = "onmouseout"
	OnMouseOver    EventType = "onmouseover"
	OnMouseUp      EventType = "onmouseup"
	OnWheel        EventType = "onwheel"
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

// AddScript adds a script by name that is triggered when a specific event occurs.
func (e *Events) AddScript(etype EventType, scriptName string) *Events {
	e.events = append(e.events, event{string(etype), scriptName})
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
