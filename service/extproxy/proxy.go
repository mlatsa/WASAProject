package extproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func MustProxy(prefix string) http.Handler {
	up := os.Getenv("UPSTREAM_API")
	u, err := url.Parse(up)
	if err != nil || up == "" {
		log.Fatalf("UPSTREAM_API is not a valid URL (got %q)", up)
	}
	p := httputil.NewSingleHostReverseProxy(u)
	p.Director = func(r *http.Request) {
		r.URL.Scheme = u.Scheme
		r.URL.Host = u.Host
		trim := strings.TrimPrefix(r.URL.Path, prefix)
		if !strings.HasPrefix(trim, "/") {
			trim = "/" + trim
		}
		r.URL.Path = trim
		r.URL.RawPath = trim
		r.Host = u.Host
		if tok := os.Getenv("UPSTREAM_TOKEN"); tok != "" {
			r.Header.Set("Authorization", "Bearer "+tok)
		}
	}
	p.ModifyResponse = func(resp *http.Response) error {
		h := resp.Header
		h.Del("Access-Control-Allow-Origin")
		h.Del("Access-Control-Allow-Methods")
		h.Del("Access-Control-Allow-Headers")
		h.Del("Access-Control-Expose-Headers")
		h.Del("Access-Control-Allow-Credentials")
		return nil
	}
	return p
}
