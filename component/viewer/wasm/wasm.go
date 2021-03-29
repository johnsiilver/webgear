package wasm

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/johnsiilver/webgear/handlers"
	wasmHTTP "github.com/johnsiilver/webgear/wasm/http"
)

// Viewer provides an HTTP server that will run to view an individual component that
// is inside a wasm binary.
type Viewer struct {
	port       int
	binName    string
	useModules bool

	mu sync.Mutex

	serveFrom *serveFrom
	h         *handlers.Mux
}

// Option provides an optional argument to New().
type Option func(v *Viewer)

type serveFrom struct {
	from string
	exts []string
}

// ServeOtherFiles looks at path "from" and serves files below that directory with the extensions in "exts".
// Extensions should be like ".png" or ".css".
func ServeOtherFiles(from string, exts []string) Option {
	return func(v *Viewer) {
		v.serveFrom = &serveFrom{from, append(exts, ".wasm")}
	}
}

// UseModules indicates if we are compiling Go using Go modules. By default this is false.
func UseModules() Option {
	return func(v *Viewer) {
		v.useModules = true
	}
}

// New constructs a new Viewer. binName must be a .wasm file that has its main package
// in the director ./main.  This will be built on startup. This can be rebuilt and
// relaunched by simply hitting 'r' in the terminal.
func New(port int, binName string, options ...Option) *Viewer {
	v := &Viewer{
		port:    port,
		binName: binName,
		h:       handlers.New(handlers.DoNotCache()),
	}

	for _, o := range options {
		o(v)
	}

	if err := v.build(); err != nil {
		panic(err)
	}

	u, err := url.Parse("/static/viewer/main/" + binName)
	if err != nil {
		panic(err)
	}

	binHandle, err := wasmHTTP.Handler(u)
	if err != nil {
		panic(err)
	}

	v.h.HTTPHandler("/", binHandle)
	if v.serveFrom != nil {
		v.h.ServeFilesFrom(v.serveFrom.from, "", v.serveFrom.exts)
	}

	return v
}

func (v *Viewer) build() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	log.Println("starting wasm rebuild")

	if v.useModules {
		os.Setenv("GO111MODULE", "on")
	} else {
		os.Setenv("GO111MODULE", "off")
	}
	os.Setenv("GOOS", "js")
	os.Setenv("GOARCH", "wasm")
	out, err := exec.Command("/usr/local/go/bin/go", "build", "-o", filepath.Join("./main", v.binName), "./main").CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Printf("Error building wasm binary: %s", err)
		for _, entry := range os.Environ() {
			log.Println(entry)
		}
		return err
	}
	return nil
}

// Run runs the viewer and will block forever.
func (v *Viewer) Run() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	go func() {
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			if string(b) == "r" {
				if err := v.build(); err == nil {
					fmt.Println("wasm rebuild SUCCEEDED, reload web browser")
				} else {
					fmt.Println("wasm rebuild FAILED: ", err)
				}
			} else {
				fmt.Printf("I don't know what %q is\n", string(b))
			}
		}
	}()

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
