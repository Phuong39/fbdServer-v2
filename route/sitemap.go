package route

import (
	"bytes"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/theTardigrade/fbdServer-v2/environment"
	"github.com/theTardigrade/fbdServer-v2/options"
)

const (
	sitemapPathsPerSub = 1_000
)

var (
	sitemapPathsMap   = make(map[string]struct{})
	sitemapPathsSlice []string
	sitemapPathsMutex sync.RWMutex
)

func sitemapPathAdd(path string) {
	defer sitemapPathsMutex.Unlock()
	sitemapPathsMutex.Lock()

	if _, found := sitemapPathsMap[path]; !found {
		sitemapPathsMap[path] = struct{}{}
		sitemapPathsSlice = append(sitemapPathsSlice, path)
	}
}

func sitemapPathAddMany(paths []string) {
	defer sitemapPathsMutex.Unlock()
	sitemapPathsMutex.Lock()

	for _, path := range paths {
		if _, found := sitemapPathsMap[path]; !found {
			sitemapPathsMap[path] = struct{}{}
			sitemapPathsSlice = append(sitemapPathsSlice, path)
		}
	}
}

func sitemapPathCount() int {
	defer sitemapPathsMutex.RUnlock()
	sitemapPathsMutex.RLock()

	return len(sitemapPathsSlice)
}

var (
	sitemapGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		siteDomain := environment.Data.MustGet("site_domain")
		subCount := int(math.Ceil(float64(sitemapPathCount()) / float64(sitemapPathsPerSub)))

		var buffer bytes.Buffer

		buffer.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
		buffer.WriteString(`<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

		for i := 1; i <= subCount; i++ {
			buffer.WriteString(`<sitemap>`)
			buffer.WriteString(`<loc>https://` + siteDomain + `/sitemap/` + strconv.Itoa(i) + `/sub.xml</loc>`)
			buffer.WriteString(`</sitemap>`)
		}

		buffer.WriteString(`</sitemapindex>`)

		w.WriteHeader(http.StatusOK)
		w.Write(buffer.Bytes())
	})
)

func init() {
	options.Options.Routes.Get["/sitemap.xml"] = sitemapGetHandler
}
