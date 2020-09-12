package rest

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	server "github.com/johnsiilver/webgear/wasm/http/rest/testdata/grpc"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/johnsiilver/webgear/wasm/http/rest/testdata/grpc/proto"
	//"github.com/kylelemons/godebug/pretty"
)

const restPort = 8080
const grpcPort = 8081

// setupGRPC sets up a GRPC service listening on addr.
func setupGRPC(addr string, dataPath string) {
	if err := os.MkdirAll(dataPath, 0700); err != nil {
		panic(err)
	}

	// Here we setup gRPC. It runs privately on localhost because only our gateway is going to dial it.
	// You could expose both if you wanted to.
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	snippetsServ := server.NewService(dataPath)
	pb.RegisterSnippetsServer(grpcServer, snippetsServ)

	go func() {
		err := grpcServer.Serve(lis)
		panic(err.Error())
	}()
}

// setupREST sets up a REST service on lis that proxies requests to our GRPC service at grpcAddr.
func setupREST(ctx context.Context, grpcAddr string) *runtime.ServeMux {
	// Now we are going to setup our reverse proxy.
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()} // Only local things need to talk to grpc, so no TLS.
	if err := pb.RegisterSnippetsHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		panic(err)
	}
	return mux
}

func addr(port int) string {
	return fmt.Sprintf("0:%d", port)
}

func init() {
	dataPath := filepath.Join(os.TempDir(), uuid.New().String())

	setupGRPC(addr(grpcPort), dataPath)
	restMux := setupREST(context.Background(), addr(grpcPort))

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", restPort),
		Handler:        handlerWrapper(restMux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("rest server serving on :%d", restPort)

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	time.Sleep(1 * time.Second) // Yeah, I know, this is sloppy
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type deflateResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w deflateResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// handlerWrapper takes a handler and wraps it in a general handler that will reject
// calls that don't have context-type set to "application/grpc-gateway" and will compress
// responses with accept-encoding "gzip" or "deflate".
func handlerWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Header.Get("content-type") {
			case "application/grpc-gateway":
			default:
				http.NotFoundHandler().ServeHTTP(w, r)
				return
			}

			// We accept incoming content compressed with gzip.
			switch r.Header.Get("Content-Encoding") {
			case "gzip":
				r.Body = ioutil.NopCloser(gzipDecompress(r.Body))
			}

			// We can compress back to the user in gzip or deflate.
			switch r.Header.Get("Accept-Encoding") {
			case "gzip":
				w.Header().Set("Content-Encoding", "gzip")
				gz := gzip.NewWriter(w)
				defer gz.Close()
				grw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
				next.ServeHTTP(grw, r)
			case "deflate":
				w.Header().Set("Content-Encoding", "deflate")
				de, _ := flate.NewWriter(w, 3)
				defer de.Close()
				drw := deflateResponseWriter{Writer: de, ResponseWriter: w}
				next.ServeHTTP(drw, r)
			default:
				next.ServeHTTP(w, r)
			}

		},
	)
}

func TestCalls(t *testing.T) {
	unixNano := time.Now().UnixNano()
	content := "hello world"
	save := "/v1/snippetsService/save"
	get := "/v1/snippetsService/get"

	headerWithGzipSend := DefaultHeaders()
	headerWithGzipSend.Add("Content-Encoding", "gzip")

	tests := []struct {
		desc    string
		path    string
		req     proto.Message
		got     proto.Message
		want    proto.Message
		options []Option
	}{
		{
			desc: "No send compression, default settings",
			path: save,
			req: &pb.SaveReq{
				UnixNano: unixNano,
				Content:  content,
			},
			got:  &pb.SaveResp{},
			want: &pb.SaveResp{},
		},
		{
			desc: "Gzip compression on send, default settings",
			path: get,
			req: &pb.GetReq{
				UnixNano: unixNano,
			},
			got: &pb.GetResp{},
			want: &pb.GetResp{
				UnixNano: unixNano,
				Content:  content,
			},
			options: []Option{
				CompressRequests(GzipCompress),
			},
		},
	}

	u, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", restPort))
	if err != nil {
		panic(err)
	}

	for _, test := range tests {
		log.Println("Doing test: ", test.desc)
		client := New(u, test.options...)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := client.Call(ctx, test.path, test.req, test.got)
		cancel()
		if err != nil {
			t.Errorf("TestCalls(%s): got err == %s, want err == nil", test.desc, err)
			continue
		}

		if diff := cmp.Diff(test.want, test.got, protocmp.Transform()); diff != "" {
			t.Errorf("TestCalls(%s): -want/+got:\n%s", test.desc, diff)
		}
	}
}
