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
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		storeName := bone.GetValue(r, "storeName")
		hashedGUID := bone.GetValue(r, "itemHashedGUID")

		item, itemFound, err := model.ItemFromHashedGUID(hashedGUID)
		if err != nil {
			panic(err)
		}
		if !itemFound || item.StoreName != storeName {
			notFoundHandler(w, r)
			return
		}

		data := dataDefault()

		data["item"] = item
		data["page_title"] = html.UnescapeString(string(item.Title)) + " | " + data["site_title"].(string)

		err = template.Views.ExecuteTemplate(w, "item", "main", data)
		if err != nil {
			panic(err)
		}
	})
)

const (
	itemPath = "/store/:storeName/item/:itemHashedGUID"
)

func init() {
	options.Options.Routes.Get[itemPath] = itemGetHandler
}
