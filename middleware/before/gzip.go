package before

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/theTardigrade/fbdServer-v2/environment"
	"github.com/theTardigrade/fbdServer-v2/options"

	grw "github.com/theTardigrade/gzipResponseWriter"
)

var (
	gzipEnableGzip = environment.Data.LazyGetBool("enable_gzip")
)

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// never use gzip unless explicitly enabled by option
		if !gzipEnableGzip {
			goto skipLabel
		}

		// check that browser accepts gzip compression
		for _, headerValue := range r.Header["Accept-Encoding"] {
			var value bool

			for _, encoding := range strings.Split(headerValue, ",") {
				if strings.TrimSpace(encoding) == "gzip" {
					value = true
					break
				}
			}

			if !value {
				goto skipLabel
			}
		}

		// never gzip for bots
		if ua := strings.ToLower(r.UserAgent()); strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") {
			goto skipLabel
		}

		// never gzip websocket connections
		for _, headerValue := range r.Header["Upgrade"] {
			if strings.Contains(headerValue, "websocket") {
				goto skipLabel
			}
		}

		// only ever compress extensionless paths or those with the given extensions
		if ext := filepath.Ext(r.URL.Path); len(ext) > 0 {
			switch ext[1:] {
			case "css", "js", "txt", "xml":
				// do nothing
			default:
				goto skipLabel
			}
		}

		// create new ResponseWriter with gzip functionality
		{
			gw, err := grw.New(w)
			if err != nil {
				panic(err)
			}
			defer gw.Close()

			gw.SetHeaders()

			w = gw
		}

	skipLabel:
		next.ServeHTTP(w, r)
	})
}

func init() {
	options.Options.Middleware.Before = append(options.Options.Middleware.Before, gzipMiddleware)
}
