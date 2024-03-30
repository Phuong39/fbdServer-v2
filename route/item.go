package route

import (
	"html"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
	"github.com/theTardigrade/fbdServer-v2/template"
)

var (
	itemGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{}

		// storeName := bone.GetValue(r, "storeName")
		hashedGUID := bone.GetValue(r, "itemHashedGUID")

		item, err := model.ItemFromHashedGUID(hashedGUID)
		if err != nil {
			panic(err)
		}

		data["item"] = item
		data["title"] = html.UnescapeString(string(item.Title)) + " | " + "Find Beautiful Designs"

		err = template.Views.ExecuteTemplate(w, "item", "main", data)
		if err != nil {
			panic(err)
		}

		// var output bytes.Buffer

		// output.WriteString(`<!DOCTYPE html><html><head><title>FBD</title>`)
		// output.WriteString(`<link rel="stylesheet" href="/static/styles/main.css" /></head>`)
		// output.WriteString(`<body>`)
		// output.WriteString(`<div class="item_profile">`)
		// output.WriteString(`<h1>` + string(item.Title) + `</h1>`)
		// output.WriteString(`<a href="` + item.LinkURL + `">`)
		// output.WriteString(`<img src=` + item.ImageURL + `" alt="` + string(item.Title) + `" />`)
		// output.WriteString(`</a>`)
		// output.WriteString(`<p>` + template.HTMLEscaper(string(item.Description)) + `</p>`)
		// output.WriteString(`</div>`)
		// output.WriteString(`</body>`)

		// w.Header().Set("Content-Type", "text/html")
		// w.WriteHeader(200)

		// w.Write(output.Bytes())
	})
)

const (
	itemPath = "/store/:storeName/item/:itemHashedGUID"
)

func init() {
	options.Options.Routes.Get[itemPath] = itemGetHandler
}
