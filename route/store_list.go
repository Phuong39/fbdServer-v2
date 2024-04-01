package route

import (
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
	"github.com/theTardigrade/fbdServer-v2/template"
)

var (
	storeListGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		data := dataDefault()

		data["page_title"] = "Store List | " + data["site_title"].(string)

		stores, err := model.StoresAll()
		if err != nil {
			panic(err)
		}

		data["stores"] = stores

		err = template.Views.ExecuteTemplate(w, "store_list", "main", data)
		if err != nil {
			panic(err)
		}
	})
)

const (
	storeListPath = "/store/list"
)

func init() {
	options.Options.Routes.Get[storeListPath] = storeListGetHandler

	sitemapPathAdd(storeListPath)
}
