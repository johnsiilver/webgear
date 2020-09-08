package content

import (
	//"context"
	"time"
	"fmt"
	"log"

	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html/builder"
	"github.com/johnsiilver/webgear/wasm"
	//"github.com/johnsiilver/webgear/wasm/examples/snippets/apps/snippets/components/content/internal/data"
	"github.com/johnsiilver/webgear/wasm/examples/snippets/grpc/date"

	. "github.com/johnsiilver/webgear/html"

)

// mode indicates if we are in view or edit mode.
type Mode int
const (
	UnknownMode = 0
	View Mode = 1
	Edit Mode = 2
)

// buildContent is used to build our content area that we use to edit and view
// our markdown content that represents notes for the week.
type buildContent struct {
	mode Mode
	immutable bool
	date time.Time
	content string

	build *builder.HTML
}

// doc is called to execute the buildContext chain and render a Doc object.
func (b *buildContent) doc() *Doc {
	b.build = builder.NewHTML(&Head{}, &Body{})
	b.build.Add(
		&Link{
			Rel: "stylesheet", 
			Href: URLParse("/static/apps/snippets/components/content/content.css"),
		},
	)

	return b.header()
}

// header assembles our header aread inside our containing div. If we are in
// edit mode we will output that we are editing the week of <date> and if viewing
// we weill leave out the "editing" part.
func (b *buildContent) header() *Doc {
	weekStr := date.WeekOf(b.date).Format("Jan 2, 2006")

	b.build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "textBanner"}})
	switch b.mode {
	case UnknownMode:
		panic("context.buildContent(): bug: received mode 'unknownMode'")
	case View:
		b.build.Into(&H{Level: 1})
		b.build.Add(TextElement(fmt.Sprintf("The Week of %s", weekStr)))
		b.build.Up()
	case Edit:
		b.build.Into(&H{Level: 1})
		b.build.Add(TextElement(fmt.Sprintf("Editing The Week of %s", weekStr)))
		b.build.Up()
	}
	b.build.Up() // Out of the div
	return b.leftArrow()
}

// leftArrow adds a "<" that can be clicked in order to go back one week.
func (b *buildContent) leftArrow() *Doc {
	log.Println("leftArrow")
	b.build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "textControls"}})

	// <
	b.build.Into(&Span{GlobalAttrs: GlobalAttrs{Class: "arrows", ID: "leftArrow"}})
	// <h1 style="display: inline">&#171</h1>
	b.build.Into(&H{Level: 1})
	b.build.Add(TextElement("&#171"))
	b.build.Up()
	b.build.Up()// Out of span

	switch b.mode {
	case View: 
		return b.viewBox()
	case Edit:
		return b.editBox()
	default:
		panic(fmt.Sprintf("unknown mode: %v", b.mode))
	}
}

// viewBox assembles our textarea tag that allows the user to view the
// content of their weekly entry. This is only called if our mode==view.
func (b *buildContent) viewBox() *Doc {
	placeholder := "Nothing to show" 
	if b.immutable {
		placeholder = "No entry for this week"
	}


	b.build.Add(
		&TextArea{
			GlobalAttrs: GlobalAttrs{ID: "textBox"},
			ReadOnly: true,
			Placeholder: placeholder,
			Element: TextElement(b.content),
		},
	)

	if !b.immutable {
		b.build.Add(
			&Input{
				GlobalAttrs: GlobalAttrs{ID: "editButton"}, 
				Type: ButtonInput, 
				Value: "Edit",
			},
		)
	}
	return b.rightArrow()
}

// editBox assembles our textarea tag that allows the user to edit the
// content of their weekly entry. This is only called if our mode==edit.
func (b *buildContent) editBox() *Doc {
	b.build.Add(
		&TextArea{
			GlobalAttrs: GlobalAttrs{ID: "textBox"},
			Name: "body",
			Element: TextElement(b.content),
		},
	)
	b.build.Up() // Out of textInput
	b.build.Up() // Out of textBox
	b.build.Add(
		&Input{
			GlobalAttrs: GlobalAttrs{ID: "editButton"}, 
			Type: ButtonInput, 
			Value: "Save",
		},
	)

	return b.rightArrow()
}

// rightArrow adds a ">" that can be clicked in order to go advance one week.
func (b *buildContent) rightArrow() *Doc {
	// >
	b.build.Into(&Span{GlobalAttrs: GlobalAttrs{Class: "arrows", ID: "leftArrow"}})
	b.build.Into(&H{Level: 1})
	b.build.Add(TextElement("&#187"))
	b.build.Up()
	b.build.Up()// Out of span

	return b.build.Doc()
}

type Args struct {
	Day time.Time
	Mode Mode
	Immutable bool
	RestEndpoint string
}

func (a Args) validate() error {
	if a.Mode == UnknownMode {
		return fmt.Errorf("Args.Mode cannot be set to UnknownMode")
	}
	if a.RestEndpoint == "" {
		return fmt.Errorf("Args.RestEndpoint cannot be set to empty string")
	}
	return nil
}

// New constructs a new component that shows a textarea ands some controls 
// that can be in a view mode or in a edit mode. The edit mode allows saving 
// to our backend server via a REST call.
func New(name string, args Args, w *wasm.Wasm, options ...component.Option) (*component.Gear, error) {
	if err := args.validate(); err != nil {
		return nil, err
	}
	if args.Day.IsZero() {
		args.Day = time.Now()
	}

	/*
	snip := data.NewSnippet(args.RestEndpoint)
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	log.Println("before fetch")
	resp, err := snip.Fetch(ctx, args.Day)
	if err != nil {
		return nil, err
	}
	*/
	log.Println("after fetch")

	content := buildContent {
		mode: args.Mode,
		immutable: args.Immutable,
		date: args.Day,
		content: "",//resp.Content,
		build: builder.NewHTML(&Head{}, &Body{}),
	}
	log.Println("before component.New()")

	gear, err := component.New(name, content.doc())
	if err != nil {
		return nil, err
	}

	return gear, nil
}