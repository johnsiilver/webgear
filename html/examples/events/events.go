package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/johnsiilver/webgear/handlers"
	"github.com/johnsiilver/webgear/html/builder"

	. "github.com/johnsiilver/webgear/html"
)

var (
	dev  = flag.Bool("dev", false, "Prevents the browser from caching content when doing development")
	port = flag.Int("port", 9568, "The port to server on")
)

func modal() []Element {
	build := builder.NewHTML(&Head{}, &Body{})

	build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "container", Class: "container"}})
	build.Into(&Div{GlobalAttrs: GlobalAttrs{Class: "card"}})
	build.Into(&Div{GlobalAttrs: GlobalAttrs{Class: "content-wrapper"}})
	build.Add(&P{Elements: []Element{TextElement("Demo Modal")}})
	build.Up()
	build.Add(
		&Span{
			GlobalAttrs: GlobalAttrs{Class: "button button__link"},
			Elements:    []Element{TextElement("Close")},
			// Here we create an event object and add a script (all inline).
			// Note that when quoting strings in the javascript, single quotes work better.
			Events: (&Events{}).AddScript(
				OnClick,
				`document.getElementById('container').style.visibility = 'hidden';`,
			),
		},
	)

	return build.Doc().Body.Elements
}

func main() {
	flag.Parse()

	doc := &Doc{
		Head: &Head{
			Elements: []Element{
				&Meta{Charset: "utf-8"},
				&Title{TagValue: TextElement("Hello World")},
				&Link{
					Rel:  "stylesheet",
					Href: URLParse("/static/index/index.css"),
				},
			},
		},
		Body: &Body{
			Elements: modal(),
		},
	}

	opts := []handlers.Option{}
	if *dev {
		opts = append(
			opts,
			handlers.DoNotCache(),
		)
	}

	h := handlers.New(opts...)

	// Serves up files ending with .css from /static/...
	h.ServeFilesWorkingDir([]string{".css"})

	// Our doc will now be served at the index page.
	h.MustHandle("/", doc)

	// Serve the content using the http.Server.
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        h.ServerMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("http server serving on :%d", *port)
	log.Fatal(server.ListenAndServe())
}
