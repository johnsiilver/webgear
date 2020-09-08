package main

import (
	"flag"

	"github.com/johnsiilver/webgear/component/viewer/wasm"
)

var (
	port    = flag.Int("port", 8080, "The port to run the server on")
)

func main() {
	flag.Parse()

	v := wasm.New(
		*port,
		"main.wasm",
		wasm.UseModules(),
		wasm.ServeOtherFiles("../", []string{".css", ".jpg", ".svg", ".png"}),
	)

	v.Run()
}
