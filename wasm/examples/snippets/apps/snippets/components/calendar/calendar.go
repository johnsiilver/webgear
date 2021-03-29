package calendar

import (
	"fmt"
	"time"

	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html/builder"
	"github.com/johnsiilver/webgear/wasm"

	. "github.com/johnsiilver/webgear/html"
)

// buildCalendar builds our HTML calendar layout.
type buildCalendar struct {
	args *Args

	build *builder.HTML
}

func (b *buildCalendar) doc() *Doc {
	b.build = builder.NewHTML(&Head{}, &Body{})

	// Add our outer stylesheet and containing div.
	b.build.Add(
		&Link{
			Rel:  "stylesheet",
			Href: URLParse(b.args.cssPath()),
		},
	)
	b.build.Into(
		&Div{
			GlobalAttrs: GlobalAttrs{ID: "calendarDiv", Class: "jzdbox1 jzdbasf jzdcal"},
		},
	)

	return b.headerBar()
}

func (b *buildCalendar) headerBar() *Doc {
	b.build.Into(&Ul{})

	if b.args.MonthArrows {
		// <
		b.build.Into(
			&Li{GlobalAttrs: GlobalAttrs{Class: "prev"}},
		)
		b.build.Add(TextElement("&#10094;"))
		b.build.Up()

		// >
		b.build.Into(
			&Li{GlobalAttrs: GlobalAttrs{Class: "next"}},
		)
		b.build.Add(TextElement("&#10095;"))
		b.build.Up()
	}

	b.build.Up() // out of the Ul

	// Month and Year.
	b.build.Into(&Div{GlobalAttrs: GlobalAttrs{Class: "jzdcalt"}})
	b.build.Add(TextElement(fmt.Sprintf("%s %d", b.args.month(), b.args.year())))
	b.build.Up()

	return b.daysOfTheWeek()
}

func (b *buildCalendar) daysOfTheWeek() *Doc {
	b.build.Add(&Span{Elements: []Element{TextElement("Su")}})
	b.build.Add(&Span{Elements: []Element{TextElement("Mo")}})
	b.build.Add(&Span{Elements: []Element{TextElement("Tu")}})
	b.build.Add(&Span{Elements: []Element{TextElement("We")}})
	b.build.Add(&Span{Elements: []Element{TextElement("Th")}})
	b.build.Add(&Span{Elements: []Element{TextElement("Fr")}})
	b.build.Add(&Span{Elements: []Element{TextElement("Sa")}})

	return b.listDays()
}

func (b *buildCalendar) listDays() *Doc {
	date := time.Date(b.args.year(), b.args.month(), 1, 0, 0, 0, 0, time.FixedZone(time.Now().Zone()))
	start := date

	// If the beginning of the month isn't a Sunday, go backwards from this day until we
	// find the first Sunday.
	for {
		if start.Weekday() != time.Sunday {
			start = start.Add(-24 * time.Hour)
			continue
		}
		break
	}

	// From that first sunday, start displaying until run out of days of the month.
	for ; start.Month() <= date.Month(); start = start.Add(24 * time.Hour) {
		dayStr := fmt.Sprintf("%d", start.Day())

		// We print out the few days in the calendar from the previous month and
		// make them no lcoick.
		if start.Month() != date.Month() {
			// Makes it unclickable and blank.
			b.build.Into(&Span{GlobalAttrs: GlobalAttrs{Class: "noclick"}})
			b.build.Add(TextElement(fmt.Sprintf("%d", start.Day())))
			b.build.Up()
			continue

		}
		// Retrieve all our day information using their DayFunc.
		day := Day{}
		if b.args.DayFunc != nil {
			day = b.args.DayFunc(start.Month(), start.Day(), start.Year())
		}

		// Puts a circle since this is a day we selected.
		var span *Span
		if day.IsSelected {
			span = &Span{
				GlobalAttrs: GlobalAttrs{ID: "day" + dayStr, Class: "circle"},
			}
		} else { // No circle
			span = &Span{GlobalAttrs: GlobalAttrs{ID: "day" + dayStr}}
		}

		// If they wanted the day to be clickable with an event.
		if day.OnClick != nil {
			b.build.Into(
				wasm.AttachListener(
					LTClick,
					false,
					day.OnClick,
					DayEventArgs{Month: start.Month(), Day: start.Day(), Year: start.Year()},
					span,
				),
			)
		} else {
			b.build.Into(span)
		}
		b.build.Add(TextElement(dayStr))
		b.build.Up()
	}
	return b.build.Doc()
}

// Args are arguments about how to setup our calendar.
type Args struct {
	// CSSPath is the path to location our css file for the calendar. If left empty
	// this will default to "/static/components/calendar/calendar.css", which is likely
	// incorrect for your application.
	CSSPath string

	// Month is the month the calendar represents.
	Month time.Month

	// Year is the year the calendar represents.
	Year int

	// DayFunc is called while building the calendar for a month. The function
	// here returns information that is used to determines how a day is displayed
	// and what
	DayFunc func(month time.Month, day, year int) Day

	// MonthArrows indciate that you want arrows that move the calendar one
	// month forward or one month backwards.
	MonthArrows bool

	// FutureTime indicates that we can display times ahead of the current time.
	FutureMonths bool
}

func (a *Args) cssPath() string {
	if a == nil || a.CSSPath == "" {
		return "/static/components/calendar/calendar.css"
	}
	return a.CSSPath
}

func (a *Args) month() time.Month {
	if a == nil || a.Month == 0 {
		return time.Now().In(time.FixedZone(time.Now().Zone())).Month()
	}
	return a.Month
}

func (a *Args) year() int {
	if a == nil || a.Year == 0 {
		return time.Now().In(time.FixedZone(time.Now().Zone())).Year()
	}
	return a.Year
}

func (a *Args) currentDay() int {
	return time.Now().In(time.FixedZone(time.Now().Zone())).Day()
}

// Day represents things about a calendar day.
type Day struct {
	// IsSelected says this day is marked as individually selected.
	IsSelected bool
	// Highlighted indicates this day is highlighted in the current calendar time.
	// This is often used to mark that this is the current day or sometimes several
	// are highlighted to show a selected week or when events occur.
	Highlighted bool
	// OnClick is the function to call when this day is clicked. If this is not set,
	// then the Day will be non-clickable. This function will receive a DayEventArgs
	// as its "arg" argument.
	OnClick WasmFunc
}

type DayEventArgs struct {
	Month time.Month
	Day   int
	Year  int
}

/*
func isNextMonthFuture(year, month int) bool {
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())
}
*/

// New constructs a new component that shows a calendar.
func New(name string, contentName string, args *Args, w *wasm.Wasm, options ...component.Option) (*component.Gear, error) {
	cal := buildCalendar{args: args}

	gear, err := component.New(name, cal.doc())
	if err != nil {
		return nil, err
	}

	return gear, nil
}
