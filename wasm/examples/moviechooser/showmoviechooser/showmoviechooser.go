package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/johnsiilver/webgear/handlers"
	httpHandler "github.com/johnsiilver/webgear/wasm/http"
)

var (
	port = flag.Int("port", 8080, "")
)

func main() {
	flag.Parse()

	urlStr, _ := url.Parse("/static/moviechooser.wasm")

	movieChooserHandler, err := httpHandler.Handler(urlStr)
	if err != nil {
		panic(err)
	}

	h := handlers.New(handlers.DoNotCache())

	server := &http.Server{
		Addr:           fmt.Sprintf("127.0.0.1:%d", *port),
		Handler:        h.ServerMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	h.ServeFilesFrom("../", "", []string{".css", ".wasm"})
	h.HTTPHandler("/", movieChooserHandler)

	log.Printf("http server serving on :%d", *port)

	log.Fatal(server.ListenAndServe())
}
