package data

import (
	"net/http"
	"context"
	"fmt"
	"time"
	"bytes"
	"io/ioutil"

	"github.com/johnsiilver/webgear/wasm/examples/snippets/grpc/date"

	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/johnsiilver/webgear/wasm/examples/snippets/grpc/proto"
)

// Snippet interacts with a GRPC service over REST for grabbing or saving content.
type Snippet struct {
	endpoint string
	client *http.Client
	url string
}

// NewSnippet is the constructor for our proto service.
func NewSnippet(endpoint string) *Snippet {
	return &Snippet{
		endpoint: endpoint,
		url: fmt.Sprintf("http://%s", endpoint),
		client: &http.Client{},
	}
}

// Fetch fetches content for a date from the server.
func (s *Snippet) Fetch(ctx context.Context, day time.Time) (*pb.GetResp, error) {
	pbReq := &pb.GetReq{
		UnixNano: date.SafeUnixNano(day),
	}

	b, err := protojson.Marshal(pbReq)
	if err != nil {
		return nil, fmt.Errorf("Snippet.Fetch() could not marshal a request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.url, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("problem fetching content for day %v: %w", day, err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Snippet.Fetch() had REST error: %w", err)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Snippet.Fetch() had error reading response body: %w", err)
	}

	pbResp := &pb.GetResp{}
	if err := protojson.Unmarshal(b, pbResp); err != nil {
		return nil, fmt.Errorf("Snippet.Fetch() had error unmarshalling response body: %w", err)
	}
	return pbResp, nil
}

// Save saves content for a date to the server.
func (s *Snippet) Save(ctx context.Context, day time.Time, content string) error {
	pbReq := &pb.SaveReq{
		UnixNano: date.SafeUnixNano(day),
		Content: content,
	}

	b, err := protojson.Marshal(pbReq)
	if err != nil {
		return fmt.Errorf("Snippet.Save() could not marshal a request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("problem saving content for day %v: %w", day, err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("Snippet.Save() had REST error: %w", err)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Snippet.Save() had error reading response body: %w", err)
	}

	pbResp := &pb.SaveResp{}
	if err := protojson.Unmarshal(b, pbResp); err != nil {
		return fmt.Errorf("Snippet.Save() had error unmarshalling response body: %w", err)
	}
	return nil
}