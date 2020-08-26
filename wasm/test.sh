#! /bin/bash
# Note: you need NodeJs installed to provide the JS environment.

GOOS=js GOARCH=wasm go test -exec "node ./wasm_exec.js"
