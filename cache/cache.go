package cache

import (
	"bytes"
	"crypto/sha1"
	"github.com/go-chi/stampede"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CreateKey creates a package specific key for a given string
func CreateKey(u string) string {
	return urlEscape("pageCache", u)
}

func urlEscape(prefix string, u string) string {
	key := url.QueryEscape(u)
	if len(key) > 200 {
		h := sha1.New()
		_, _ = io.WriteString(h, u)
		key = string(h.Sum(nil))
	}
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	return buffer.String()
}

// Authorization
// Include anything user specific, e.g. Authorization Token
var customKeyFunc = func(r *http.Request) uint64 {
	token := r.Header.Get("Authorization")
	return stampede.StringToHash(r.Method, strings.ToLower(strings.ToLower(token)))
}

type cache struct {
	isAuthorization bool
	cacheSize       int
	ttl             time.Duration
	paths           []string
}

func NewCache(isAuthorization bool) *cache {
	return &cache{
		isAuthorization: isAuthorization,
		cacheSize:       512,
		ttl:             time.Minute * 10,
	}
}

func (c *cache) SetCacheSize(cacheSize int) {
	c.cacheSize = cacheSize
}

func (c *cache) SetTtl(ttl time.Duration) {
	c.ttl = ttl
}

func (c *cache) SetPaths(paths []string) {
	c.paths = paths
}

// CachePage strings.ToLower(r.URL.Path)
func (c cache) CachePage() func(next http.Handler) http.Handler {
	if c.isAuthorization {
		return stampede.HandlerWithKey(c.cacheSize, c.ttl, customKeyFunc)
	} else {

		return stampede.Handler(c.cacheSize, c.ttl, c.paths...)
	}
}

func (c cache) CachePageWithQuery() func(next http.Handler) http.Handler {
	if c.isAuthorization {
		// Authorization
		// Include anything user specific, e.g. Authorization Token
		fuc := func(r *http.Request) uint64 {
			token := r.Header.Get("Authorization") + r.URL.RequestURI()
			return stampede.StringToHash(r.Method, strings.ToLower(strings.ToLower(token)))
		}
		return stampede.HandlerWithKey(c.cacheSize, c.ttl, fuc)
	} else {
		fuc := func(r *http.Request) uint64 {
			token := r.URL.RequestURI()
			return stampede.StringToHash(r.Method, strings.ToLower(strings.ToLower(token)))
		}
		return stampede.HandlerWithKey(c.cacheSize, c.ttl, fuc)
	}
}
