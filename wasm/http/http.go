package http

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"log"

	"github.com/johnsiilver/webgear/html"
)

// handler implements http.Handler by serving up an *html.Doc.
type handler struct {
	doc *html.Doc
}

// ServeHTTP implements http.Handler.ServeHTTP.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("doing the wasm download")
	defer log.Println("done")
	h.doc.Execute(r.Context(), w, r)
}

// Handler will return a Handler that will return code to load your WASM application.
func Handler(downloadAppFrom *url.URL) (http.Handler, error) {
	p := downloadAppFrom.String()
	if p == "" {
		return nil, fmt.Errorf("the url passed(%s) to wasm.Handler was invalid", p)
	}

	doc := &html.Doc{
		Head: &html.Head{
			Elements: []html.Element{
				&html.Meta{Charset: "UTF-8"},
				&html.Script{
					TagValue: template.JS(
						fmt.Sprintf(
							`
%s

const go = new Go();
WebAssembly.instantiateStreaming(fetch("%s"), go.importObject).then((result) => {
	go.run(result.instance);
});
`, wasmExec, p),
					),
				},
			},
		},
		Body: &html.Body{},
	}
	if err := doc.Init(); err != nil {
		panic(err)
	}

	return handler{doc: doc}, nil
}
