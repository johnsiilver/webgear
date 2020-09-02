package main

import (
	"flag"
	"strings"

	"github.com/johnsiilver/webgear/component/viewer"
	"github.com/johnsiilver/webgear/wasm/examples/moviechooser/components/list"
)

var (
	port    = flag.Int("port", 8080, "The port to run the server on")
	filters = flag.String("filters", "", "A comma deliminated list of filters to apply")
)

func main() {
	flag.Parse()

	filterList := strings.Split(*filters, ",")
	for i, item := range filterList {
		filterList[i] = strings.TrimSpace(item)
		// I should check for filter correctness, but this is a demo.
	}

	list, err := list.New("list-component", filterList)
	if err != nil {
		panic(err)
	}

	v := viewer.New(
		8080,
		list,
		viewer.BackgroundColor("white"),
		viewer.ServeOtherFiles("../../../", []string{".css", ".jpg", ".svg", ".png"}),
	)

	v.Run()
}
