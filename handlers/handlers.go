// Package handlers provides handlers which can execute pages defined in the webgear package.
package handlers

import (
	"compress/gzip"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/johnsiilver/webgear/html"
)

// Mux allows building a http.ServeMux for use in serving html.Doc object output.
type Mux struct {
	mux *http.ServeMux

	caching   bool
	gzipFiles bool

	gzPool sync.Pool
}

// Option is an optional argument to the New constructor.
type Option func(m *Mux)

// DoNotCache tells the brower not to cache content.  This is extremely useful when you are doing development.
func DoNotCache() Option {
	return func(m *Mux) {
		m.caching = false
	}
}

// DoNotCompress prevents the muxer from gzip compressing content.
func DoNotCompress() Option {
	return func(m *Mux) {
		m.gzipFiles = false
	}
}

// New creates a new instance of Mux.
func New(options ...Option) *Mux {
	m := &Mux{
		mux:       http.NewServeMux(),
		caching:   true,
		gzipFiles: true,
		gzPool: sync.Pool{
			New: func() interface{} {
				return gzip.NewWriter(nil)
			},
		},
	}

	for _, option := range options {
		option(m)
	}

	return m
}

// ServerMux returns an http.ServerMux wrapped in various handlers.  Use this with http.Server{} to serve the content.
func (m *Mux) ServerMux() http.Handler {
	return m.preventCaching(m.gzip(m.mux))
}

// Handle registers the doc for a given pattern. If a handler already exists for pattern, Handle panics.
// All handles will be gzip compressed by default.
func (m *Mux) Handle(pattern string, doc *html.Doc) error {
	if err := doc.Init(); err != nil {
		return err
	}

	m.mux.Handle(
		pattern,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if err := doc.Execute(r.Context(), w, r); err != nil {
					log.Println(err)
					//http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			},
		),
	)
	return nil
}

// MustHandle is like Handle() except an error causes a panic. Returns the *Mux object so these can be chained.
func (m *Mux) MustHandle(pattern string, doc *html.Doc) *Mux {
	if err := m.Handle(pattern, doc); err != nil {
		panic(err)
	}

	return m
}

// HTTPHandler registers a standard http.Handler for the pattern on the http.ServeMux.
func (m *Mux) HTTPHandler(pattern string, handler http.Handler) {
	m.mux.Handle(pattern, handler)
}

// ServeFilesWorkingDir will serve all files with the following file extensions that are in the
// working directory or in any directory lower in the tree. It will never serve ., .. or .go files.
// These files are all served from pattern. All files are served out of the /static/ path.
func (m *Mux) ServeFilesWorkingDir(exts []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	allowed := make(map[string]bool, len(exts))
	for _, v := range exts {
		allowed[v] = true
	}

	m.mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(
				fileSystem{
					http.Dir(wd),
					allowed,
				},
			),
		),
	)
}

// ServeFilesFrom will serve all files with the following file extensions that are in the directory dir.
// It will never serve ., .. or .go files. All files are served out of the /{{root}}/ path. If root == "",
// /static/ will be used.
// Note: if called multiple times or used with ServeFilesWorkingDir(), if there are two directories
// within the top level directory with the same name and same root, you will get a collision that will
// cause a panic.
// Aka, if you do: ServeFilesFrom("/some_dir", "", []string{".img"}}) and
// ServeFilesFrom("/another_dir", "", []string{".img"}}), where /some_dir and /another_dir both contain img/, this
// will panic.
func (m *Mux) ServeFilesFrom(dir, root string, exts []string) {
	if root == "" {
		root = "/static/"
	}

	allowed := make(map[string]bool, len(exts))
	for _, v := range exts {
		allowed[v] = true
	}

	m.mux.Handle(
		root,
		http.StripPrefix(
			root,
			http.FileServer(
				fileSystem{
					http.Dir(dir),
					allowed,
				},
			),
		),
	)
}

func (m *Mux) gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.gzipFiles {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := m.gzPool.Get().(*gzip.Writer)
		defer m.gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		next.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// preventCaching will set the header so our responses will not be cached by the client.
func (m *Mux) preventCaching(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.caching {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}

		next.ServeHTTP(w, r)
	})
}
