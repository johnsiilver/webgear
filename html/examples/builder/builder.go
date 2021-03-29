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

func main() {
	flag.Parse()

	// Note: You can also build the head with the builder package.
	head := Head{
		Elements: []Element{
			&Meta{Charset: "utf-8"},
			&Title{TagValue: TextElement("Hello World")},
			&Link{
				Rel:  "stylesheet",
				Href: URLParse("/static/index/index.css"),
			},
		},
	}

	// Create a builder.HTML object with an initial head and body.
	build := builder.NewHTML(&head, &Body{})

	// Add a single H1 element containing some text.
	build.Add(
		&H{Level: 1, Elements: []Element{TextElement("Movies I Like")}},
	)

	// Create a div with a set ID and move the builders context inside the div.
	build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "myTableDiv"}})
	// Create a table in the div and move the builders context inside the table.
	build.Into(&Table{})
	// Add a tr to the table and move into it.
	build.Into(&TR{})
	// Add two th elements to the table.
	build.Add(
		&TH{Element: TextElement("Movie")},
		&TH{Element: TextElement("Category")},
	)
	// Use Up() to set the context back to the table.
	build.Up() // We are now inside the table
	build.Into(&TR{})
	build.Add(
		&TD{Element: TextElement("Blade Runner")},
		&TD{Element: TextElement("SciFi")},
	)
	build.Up() // We are now inside the table
	build.Into(&TR{})
	build.Add(
		&TD{Element: TextElement("Memento")},
		&TD{Element: TextElement("Drama")},
	)
	build.Up().Up() // We are now inside the div

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
	h.MustHandle("/", build.Doc())

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
