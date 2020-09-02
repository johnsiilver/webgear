// +build !js,!wasm

package html

/*
This file holds Execute() and ExecuteAsGear() for all compilation targets that are !wasm.
*/

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Execute executes the internal templates and writes the output to the io.Writer. This is thread-safe.
func (d *Doc) Execute(ctx context.Context, w io.Writer, r *http.Request) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
			return
		}
	}()

	if !d.initDone {
		return fmt.Errorf("Doc object did not have .Init() called before Execute()")
	}

	pipe := NewPipeline(ctx, r, w)
	pipe.Self = d

	if err := docTmpl.Execute(w, pipe); err != nil {
		return err
	}

	return pipe.HadError()
}

// ExecuteAsGear uses the Pipeline provided instead of creating one internally. This is for internal use only
// and no guarantees are made on its operation or that it will exist in the future. This is thread-safe.
func (d *Doc) ExecuteAsGear(pipe Pipeline) string {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if !d.initDone {
		pipe.Error(fmt.Errorf("Doc object did not have .Init() called before Execute()"))
		return EmptyString
	}
	pipe.Self = d

	if err := docTmpl.Execute(pipe.W, pipe); err != nil {
		pipe.Error(err)
	}
	return EmptyString
}
