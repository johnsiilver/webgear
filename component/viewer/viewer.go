package viewer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/handlers"
	"github.com/johnsiilver/webgear/html"
)

// Viewer provides an HTTP server that will run to view an individual component without having to spin up the
// entire page and service. This allows debugging a component individually.
type Viewer struct {
	port  int
	doc   *html.Doc
	color string

	serveFrom *serveFrom
	h         *handlers.Mux
}

// Option provides an optional argument to New().
type Option func(v *Viewer)

// UseDoc says to use this custom doc object and add the *Gear as the last element in the doc's Body.Elements list.
func UseDoc(doc *html.Doc) Option {
	return func(v *Viewer) {
		v.doc = doc
	}
}

// BackgroundColor changes the default background color from white to the color passed. This helps for
// styles utilize white and can't be seen. Does not work if UseDoc() was passed as an option.
func BackgroundColor(color string) Option {
	return func(v *Viewer) {
		v.color = color
	}
}

type serveFrom struct {
	from string
	exts []string
}

// ServeOtherFiles looks at path "from" and serves files below that directory with the extensions in "exts".
// Extensions should be like ".png" or ".css".
func ServeOtherFiles(from string, exts []string) Option {
	return func(v *Viewer) {
		v.serveFrom = &serveFrom{from, exts}
	}
}

// dynamicColor is used to implement an html.DynamicFunc so that we can change the background color based on the
// choice of color background they want.
type dynamicColor struct {
	color string
}

func (d dynamicColor) Color(pipe html.Pipeline) []html.Element {
	if d.color != "" {
		return []html.Element{
			&html.Style{
				TagValue: template.CSS(fmt.Sprintf("body{background-color: %s;}", d.color)),
			},
		}
	}
	return nil
}

// New constructs a new Viewer.
func New(port int, gear *component.Gear, options ...Option) *Viewer {
	v := &Viewer{
		port: port,
		h:    handlers.New(handlers.DoNotCache()),
	}

	for _, o := range options {
		o(v)
	}

	if v.doc == nil {
		v.doc = &html.Doc{
			Head: &html.Head{
				Elements: []html.Element{
					&html.Meta{Charset: "UTF-8"},
				},
			},
			Body: &html.Body{
				Elements: []html.Element{
					html.Dynamic(dynamicColor{v.color}.Color),
					gear,
					&html.Component{TagType: template.HTMLAttr(gear.Name())},
				},
			},
		}
	}
	if err := v.doc.Init(); err != nil {
		panic(err)
	}

	v.h.MustHandle("/", v.doc)
	if v.serveFrom != nil {
		v.h.ServeFilesFrom(v.serveFrom.from, "", v.serveFrom.exts)
	}

	return v
}

// Run runs the viewer and will block forever.
func (v *Viewer) Run() {
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", v.port),
		Handler:        v.h.ServerMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("http server serving on :%d", v.port)

	log.Fatal(server.ListenAndServe())
}
