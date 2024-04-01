package before

import (
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/options"
)

func stripSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		pathLen := len(path)
		var i int

	escapeLabel:
		for i = pathLen - 1; i > 0; i-- {
			switch path[i] {
			case '/', '\\':
				// noop
			default:
				break escapeLabel
			}
		}

		if i++; i == pathLen {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, path[:i], http.StatusMovedPermanently)
		}
	})
}

func init() {
	options.Options.Middleware.Before = append(options.Options.Middleware.Before, stripSlashMiddleware)
}
