package grpc

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/johnsiilver/webgear/wasm/examples/snippets/grpc/date"

	pb "github.com/johnsiilver/webgear/wasm/examples/snippets/grpc/proto"
)

// Service provides a gRPC service that can save and fetch content representing a weekly entry
// of snippets.
type Service struct {
	rootStorage string
}

// NewService is the constructor for Service.
func NewService(storagePath string) *Service {
	return &Service{
		rootStorage: storagePath,
	}
}

// Get gets content for a specific week that is saved in the filesystem.
func (s *Service) Get(ctx context.Context, in *pb.GetReq) (*pb.GetResp, error) {
	utime := date.SafeUnixNano(time.Unix(0, in.UnixNano))
	fp := filepath.Join(s.rootStorage, strconv.Itoa(int(utime)))
	_, err := os.Stat(fp)
	if err != nil {
		return &pb.GetResp{UnixNano: utime}, nil
	}
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		err = fmt.Errorf("content for file %s exists, but could not access it: %w", err)
		log.Println(err)
		return nil, err
	}
	return &pb.GetResp{UnixNano: utime, Content: string(b)}, nil
}

// Save saves content for a specific week in the filesystem.
func (s *Service) Save(ctx context.Context, in *pb.SaveReq) (*pb.SaveResp, error) {
	utime := date.SafeUnixNano(time.Unix(0, in.UnixNano))
	fp := filepath.Join(s.rootStorage, strconv.Itoa(int(utime)))
	err := ioutil.WriteFile(fp, []byte(in.Content), 0700)
	if err != nil {
		err = fmt.Errorf("we could not write content for entry %s due to disk error: %w", date.WeekOf(time.Unix(0, in.UnixNano)), err)
		log.Println(err)
		return nil, err
	}
	return &pb.SaveResp{}, nil
}
