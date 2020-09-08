package main

import (
	"fmt"
	"context"
	"log"
	"time"
	"syscall/js"

	"github.com/johnsiilver/webgear/wasm"
	"github.com/johnsiilver/webgear/wasm/examples/snippets/apps/snippets/components/calendar"

	. "github.com/johnsiilver/webgear/html"
)

func exampleDay(month time.Month, day, year int) calendar.Day {
	calDay := calendar.Day{}

	now := time.Now()
	if year == now.Year() && month  == now.Month() && day == now.Day() {
		calDay.IsSelected = true
		calDay.OnClick = func(this js.Value, root js.Value, args interface{}) {
			ourArgs := args.(calendar.DayEventArgs)
			js.Global().Call("alert", fmt.Sprintf("It is: %s %d %d", ourArgs.Month, ourArgs.Day, ourArgs.Year))
		}
	}
	return calDay
}

func main() {
	w := wasm.New()

	calendarGear, err := calendar.New(
		"calendar-component", 
		"content-component", 
		&calendar.Args{
			CSSPath: "/static/calendar.css",
			DayFunc: exampleDay,
		}, 
		w,
	)
	if err != nil {
		log.Println(err)
		return
	}

	doc := &Doc{
		Head: &Head{
			Elements: []Element{
				&Meta{Charset: "UTF-8"},
			},
		},
		Body: &Body{
			Elements: []Element{
				calendarGear,
				&Component{Gear: calendarGear},
			},
		},
	}
	w.SetDoc(doc)

	w.Run(context.Background())
}
