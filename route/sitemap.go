package route

import (
	"bytes"
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/environment"
	"github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
)

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

		buffer.WriteString(`<url>`)
		buffer.WriteString(`<loc>https://` + environment.Data.MustGet("site_domain") + `/stores</loc>`)
		buffer.WriteString(`</url>`)

		items, err := model.ItemMultipleAtRandom(2_000)
		if err != nil {
			panic(err)
		}

		for _, item := range items {
			buffer.WriteString(`<url>`)
			buffer.WriteString(`<loc>https://` + environment.Data.MustGet("site_domain") + `/store/` + item.StoreName + `/item/` + item.HashedGUID + `</loc>`)
			buffer.WriteString(`</url>`)
		}

		buffer.WriteString(`</urlset>`)

		w.WriteHeader(http.StatusFound)
		w.Write(buffer.Bytes())
	})
)

func init() {
	options.Options.Routes.Get["/sitemap.xml"] = sitemapGetHandler
}
