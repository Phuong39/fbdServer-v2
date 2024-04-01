package route

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/theTardigrade/fbdServer-v2/environment"
	"github.com/theTardigrade/fbdServer-v2/options"
)

var (
	sitemapPaths      = make(map[string]struct{})
	sitemapPathsMutex sync.RWMutex
)

func sitemapPathAdd(url string) {
	defer sitemapPathsMutex.Unlock()
	sitemapPathsMutex.Lock()

	sitemapPaths[url] = struct{}{}
}

func sitemapPathAddMany(urls []string) {
	defer sitemapPathsMutex.Unlock()
	sitemapPathsMutex.Lock()

	for _, u := range urls {
		sitemapPaths[u] = struct{}{}
	}
}

func sitemapPathCount() int {
	defer sitemapPathsMutex.RUnlock()
	sitemapPathsMutex.RLock()

	return len(sitemapPaths)
}

var (
	sitemapGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		var buffer bytes.Buffer

		buffer.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
		buffer.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

		func(siteDomain string) {
			defer sitemapPathsMutex.RUnlock()
			sitemapPathsMutex.RLock()

			for path := range sitemapPaths {
				buffer.WriteString(`<url>`)
				buffer.WriteString(`<loc>https://`)
				buffer.WriteString(siteDomain)
				buffer.WriteString(path)
				buffer.WriteString(`</loc>`)
				buffer.WriteString(`</url>`)
			}
		}(environment.Data.MustGet("site_domain"))

		buffer.WriteString(`</urlset>`)

		w.WriteHeader(http.StatusOK)
		w.Write(buffer.Bytes())
	})
)

func init() {
	options.Options.Routes.Get["/sitemap.xml"] = sitemapGetHandler
}
