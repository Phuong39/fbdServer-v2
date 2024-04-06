package route

import (
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
	"github.com/theTardigrade/fbdServer-v2/template"
)

var (
	homeGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		data := dataDefault()

		data["page_title"] = data["site_title"]
		data["page_description"] = `A random selection of high-quality customizable items is listed here on ` + data["site_title"].(string) + `.`

		items, err := model.ItemMultipleAtRandom(120)
		if err != nil {
			panic(err)
		}
		data["items"] = items

		template.Views.ExecuteTemplate(w, "store_page", "main", data)

		// var output bytes.Buffer

		// output.WriteString(`<!DOCTYPE html><html><head><title>FBD</title><link rel="stylesheet" href="/static/styles/main.css" /></head>`)
		// output.WriteString(`<body>`)
		// output.WriteString(`<div class="items">`)

		// items, err := model.ItemMultipleAtRandom(240)
		// if err != nil {
		// 	panic(err)
		// }

		// for _, item := range items {
		// 	output.WriteString(`<a href="/store/` + item.StoreName + `/item/` + item.HashedGUID + `">`)
		// 	output.WriteString(`<div class="item">`)
		// 	output.WriteString(`<h1>` + string(item.Title) + `</h1>`)
		// 	output.WriteString(`<img src=` + item.ImageURL + `" alt="` + string(item.Title) + `" />`)
		// 	output.WriteString(`</div>`)
		// 	output.WriteString(`</a>`)
		// }

		// output.WriteString(`</div>`)
		// output.WriteString(`</body>`)

		// resp.Header().Set("Content-Type", "text/html")
		// resp.WriteHeader(200)

		// resp.Write(output.Bytes())
	})
)

const (
	homePath = "/"
)

func init() {
	options.Options.Routes.Get[homePath] = homeGetHandler

	sitemapPathAdd(homePath)
}
