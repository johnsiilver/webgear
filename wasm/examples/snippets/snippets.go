package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/johnsiilver/webgear/handlers"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	server "github.com/johnsiilver/webgear/wasm/examples/snippets/grpc"
	httpHandler "github.com/johnsiilver/webgear/wasm/http"

	pb "github.com/johnsiilver/webgear/wasm/examples/snippets/grpc/proto"
)

var (
	grpcPort = flag.Int("grpc_port", 9000, "The port @ 127.0.0.1 to run on grpc on")
	port     = flag.Int("port", 8080, "The port to run REST on")
	httpPort = flag.Int("http_port", 8081, "The port to run our http pages on")
)

const snippetsDir = "/tmp/snippets"

// setupGRPC sets up a GRPC service listening on addr.
func setupGRPC(addr string) {
	// Here we setup gRPC. It runs privately on localhost because only our gateway is going to dial it.
	// You could expose both if you wanted to.
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	snippetsServ := server.NewService(snippetsDir)
	pb.RegisterSnippetsServer(grpcServer, snippetsServ)

	go func() {
		err := grpcServer.Serve(lis)
		panic(err.Error())
	}()
}

// setupREST sets up a REST service on lis that proxies requests to our GRPC service at grpcAddr.
func setupREST(ctx context.Context, addr string, grpcAddr string) {
	// Now we are going to setup our reverse proxy.
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()} // Only local things need to talk to grpc, so no TLS.
	if err := pb.RegisterSnippetsHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		panic(err)
	}

	httpServ := &http.Server{Addr: addr, Handler: mux}

	go func() {
		// THIS IS INSECURE, BECAUSE THERE IS NO TLS OR AUTH.
		if err := httpServ.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}

// setupWASM registers handlers that will download our wasm app when someone goes to /.
// We could register multiple wasm apps at different URLs, but we aren't doing that.
// We are also running it on its own address:port. Running all these on the same port
// is a pain, because we have to do a lot of weird things because of http2 and not using TLS.
// At some point I'll wrap all these things up to get rid of the boiler plate.
func setupWASM(addr, restAddr string) {
	// Create an http.Handler for downloading and running our wasm app.
	urlStr, _ := url.Parse("/static/apps/snippets/snippets.wasm")
	snippetsHandler, err := httpHandler.Handler(urlStr)
	if err != nil {
		panic(err)
	}

	h := handlers.New(handlers.DoNotCache())
	server := &http.Server{
		Addr:           addr,
		Handler:        h.ServerMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Serve all .css and .wasm files.
	h.ServeFilesFrom("", "", []string{".css", ".wasm", ".svg"})
	// Serve our snippet app from /.
	h.HTTPHandler("/", snippetsHandler)

	log.Printf("http server serving on at %s", addr)

	go log.Fatal(server.ListenAndServe())
}

func main() {
	flag.Parse()
	ctx := context.Background()

	grpcAddr := fmt.Sprintf("127.0.0.1:%d", *grpcPort)
	gatewayAddr := fmt.Sprintf(":%d", *port)
	httpAddr := fmt.Sprintf(":%d", *httpPort)

	if _, err := os.Stat(snippetsDir); err != nil {
		if err := os.MkdirAll(snippetsDir, 0700); err != nil {
			panic(fmt.Sprintf("cannot create directory(%s) to store snippets in: %s", snippetsDir, err))
		}
	}

	setupGRPC(grpcAddr)
	setupREST(ctx, gatewayAddr, grpcAddr)
	setupWASM(httpAddr, gatewayAddr)

	select {}
}
