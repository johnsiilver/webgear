package banner

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/johnsiilver/go_basics/site/config"
	"github.com/johnsiilver/webgear/html"
)

func TestBanner(t *testing.T) {
	conf := &config.VideoFiles{
		&config.VideoFile{
			Name: "Video 0",
			URL:  "/video/0",
		},
	}
	req := &http.Request{
		URL: html.URLParse("/video/0"),
	}

	buff := &strings.Builder{}
	pipe := html.NewPipeline(context.Background(), req, buff)

	b, err := New("my-banner", conf)
	if err != nil {
		panic(err)
	}

	out, err := ioutil.ReadFile("testdata/want")
	if err != nil {
		panic(err)
	}

	space := regexp.MustCompile(`\s+`)

	b.Execute(pipe)

	got := strings.TrimSpace(space.ReplaceAllString(buff.String(), " "))
	want := strings.TrimSpace(space.ReplaceAllString(string(out), " "))

	if got != want {
		t.Errorf("TestBanner: want:\n%s\ngot:\n%s", want, got)
	}
}
