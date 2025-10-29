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
		r.Host = u.Host
		if tok := os.Getenv("UPSTREAM_TOKEN"); tok != "" {
			r.Header.Set("Authorization", "Bearer "+tok)
		}
	}
	return p
}
