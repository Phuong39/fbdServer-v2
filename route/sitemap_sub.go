package route

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/theTardigrade/fbdServer-v2/environment"
	"github.com/theTardigrade/fbdServer-v2/options"
	basicServer "github.com/theTardigrade/golang-basicServer"
)

var (
	sitemapSubGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		numberRaw := basicServer.ValueFromRequest(r, "number")
		if numberRaw == "" {
			notFoundHandler(w, r)
			return
		}

		number, err := strconv.Atoi(numberRaw)
		if err != nil {
			notFoundHandler(w, r)
			return
		}

		var buffer bytes.Buffer
		var urlCount int

		buffer.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
		buffer.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

		func(siteDomain string) {
			startIndex := sitemapPathsPerSub * (number - 1)
			endIndex := sitemapPathsPerSub * number

			defer sitemapPathsMutex.RUnlock()
			sitemapPathsMutex.RLock()

			if endIndex > len(sitemapPathsSlice) {
				endIndex = len(sitemapPathsSlice)
			}

			if startIndex < len(sitemapPathsSlice) {
				for i := startIndex; i < endIndex; i++ {
					path := sitemapPathsSlice[i]

					buffer.WriteString(`<url>`)
					buffer.WriteString(`<loc>https://`)
					buffer.WriteString(siteDomain)
					buffer.WriteString(path)
					buffer.WriteString(`</loc>`)
					buffer.WriteString(`</url>`)

					urlCount++
				}
			}
		}(environment.Data.MustGet("site_domain"))

		if urlCount == 0 {
			notFoundHandler(w, r)
			return
		}

		buffer.WriteString(`</urlset>`)

		w.WriteHeader(http.StatusOK)
		w.Write(buffer.Bytes())
	})
)

func init() {
	options.Options.Routes.Get["/sitemap/:number/sub.xml"] = sitemapSubGetHandler
}
