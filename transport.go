package restc

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gregjones/httpcache"
	"github.com/pkg/errors"
)

type DebugTransport struct {
	Transport           http.RoundTripper
	Cache               httpcache.Cache
	MarkCachedResponses bool
}

func NewDebugTransport(c httpcache.Cache) *DebugTransport {
	return &DebugTransport{Cache: c, MarkCachedResponses: true}
}

func generateCacheKey(req *http.Request) (string, error) {
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return "", errors.Wrap(err, "dump request")
	}
	return fmt.Sprintf("%s-%s-%x", req.Method, req.URL.String(), md5.Sum(reqBytes)), nil
}

func (t *DebugTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	cacheKey, err := generateCacheKey(req)

	if err != nil {
		return nil, errors.Wrap(err, "generate cache key")
	}

	if cachedVal, ok := t.Cache.Get(cacheKey); ok {
		b := bytes.NewBuffer(cachedVal)
		if cachedResp, err := http.ReadResponse(bufio.NewReader(b), req); err == nil {
			if t.MarkCachedResponses {
				cachedResp.Header.Set("X-From-Cache", "1")
			}
			return cachedResp, nil
		} else {
			return nil, errors.Wrap(err, "get cached response")
		}
	}

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	resp, err = transport.RoundTrip(req)
	if err != nil {
		return nil, errors.Wrap(err, "round trip")
	}

	respBytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, errors.Wrap(err, "dump response")
	}

	t.Cache.Set(cacheKey, respBytes)

	return resp, nil
}
