/*
Package handlers provides http.Handler(s) which can execute pages defined in the webgear package.  All data is
returned gzip compressed.

Usage is as follows:
	// Create new handlers.Mux object with an option to tell clients not to cache results.
	h := handlers.New(handlers.DoNotCache())

	// Serve all files from the the binary working directory and below it (recursively) that have
	// the file extensions listed.
	h.ServeFilesWorkingDir([]string{".css", ".jpg", ".svg", ".png"})

	// Create a *html.Doc object that we want to serve from "/".
	index, err := index.New(conf)
	if err != nil {
		panic(err)
	}

	// Attach that object to /.
	h.MustHandle("/", index)

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
*/
package handlers

import (
	"compress/gzip"
	"io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/johnsiilver/webgear/html"
)

// Mux allows building a http.ServeMux for use in serving html.Doc object output.
type Mux struct {
	mux *http.ServeMux

	expireCache *expireCache

	caching   bool
	gzipFiles bool
	debug     bool

	gzPool sync.Pool
}

// Option is an optional argument to the New constructor.
type Option func(m *Mux)

// StaticMode causes the server to cache the output from any page after the first call.  
// The static page in cache will be returned as long as the expire time hasn't passed
// since the last call to a page. The expire time is simply there to keep pages that aren't
// used often from costing memory. expire and sweep must >= 30 seconds.
func StaticMode(expire, sweep time.Duration) Option {
	if expire < 30 * time.Second {
		panic("expire must be >= 30 seconds")
	}
	if sweep < 30 * time.Second {
		panic("expire must be >= 30 seconds")
	}
	return func(m *Mux) {
		m.expireCache = newExpireCache(expire, sweep)
	}
}

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

// Debug causes error messages from HTML rendering to be output.
func Debug() Option {
	return func(m *Mux) {
		m.debug = true
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
	return m.staticCache(
		m.preventCaching(
			m.gzip(m.mux),
		),
	)
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
				r.ParseForm()
				if err := doc.Execute(r.Context(), w, r); err != nil {
					if m.debug {
						log.Println(err)
					}
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

// ServeFS passes a fs.FS that is walked and servers out of a root of /static/. This is similar to
// ServeFilesWorkingDir() except it serves up all files in the FS that can be walked. Generally this
// if for embeded files. Cannot be used with ServeFilesWorkingDir().
func (m *Mux) ServeFS(filesys fs.FS) {
	m.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(filesys))))
}

// ServeFilesWorkingDir will serve all files with the following file extensions that are in the
// working directory or in any directory lower in the tree. It will never serve ., .. or .go files.
// These files are all served from pattern. All files are served out of the /static/ path.
// Cannot be used with ServeFS().
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

func (m *Mux) staticCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// No cache.
		if m.expireCache == nil {
			next.ServeHTTP(w, r)
			return
		}

		// Cache hit.
		e, err := m.expireCache.get(r.URL)
		if err == nil {
			for k, v := range e.h {
				w.Header()[k] = v
			}
			w.Write(e.b)
			return
		}

		// Cache miss.
		c := httptest.NewRecorder()
		next.ServeHTTP(c, r)

		for k, v := range c.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(c.Code)
		content := c.Body.Bytes()

		m.expireCache.put(r.URL, c.HeaderMap, content)
		w.Write(content)
	})
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
