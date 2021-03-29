package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/johnsiilver/webgear/handlers"

	"github.com/johnsiilver/webgear/html/examples/components/index"
)

var (
	dev  = flag.Bool("dev", false, "Prevents the browser from caching content when doing development")
	port = flag.Int("port", 9568, "The port to server on")
)

func main() {
	flag.Parse()

	doc, err := index.New()
	if err != nil {
		panic(err)
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
	h.ServeFilesWorkingDir([]string{".svg", ".css"})

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
