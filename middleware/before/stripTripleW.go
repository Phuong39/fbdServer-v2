package before

import (
	"net/http"
	"strings"

	"github.com/theTardigrade/fbdServer-v2/options"
)

func stripTripleWMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if host := r.Host; strings.Contains(host, "www.") {
			url := r.URL
			url.Host = strings.ReplaceAll(host, "www.", "")

			http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func init() {
	options.Options.Middleware.Before = append(options.Options.Middleware.Before, stripTripleWMiddleware)
}
