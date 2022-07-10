package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
	"sync"
)

type cacheEntry struct {
	h http.Header
	b []byte
	t time.Time
}

type expireCache struct {
	expire time.Duration
	sweep time.Duration
	
	sync.Mutex
	pages map[string]cacheEntry
}

func newExpireCache(expire, sweep time.Duration) *expireCache {
	return &expireCache{pages: map[string]cacheEntry{}}
}

func (e *expireCache) get(u *url.URL) (cacheEntry, error) {
	s := u.String()

	e.Lock()
	defer e.Unlock()


	ec, ok := e.pages[s]
	if ok {
		ec.t = time.Now().Add(e.expire)
		e.pages[s] = ec
		return ec, nil
	}
	return cacheEntry{}, fmt.Errorf("not found")
}

func (e *expireCache) put(u *url.URL, h http.Header, b []byte) {
	e.Lock()
	defer e.Unlock()

	e.pages[u.String()] = cacheEntry{h: h, b: b, t: time.Now()}
}

func (e *expireCache) purge() {
	n := time.Now()
	e.Lock()
	defer e.Unlock()
	for _ = range time.Tick(e.sweep) {
		for k, ec := range e.pages {
			if ec.t.After(n) {
				delete(e.pages, k)
			}
		}
	}
}
