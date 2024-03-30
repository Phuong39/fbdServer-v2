package route

import (
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
	"github.com/theTardigrade/fbdServer-v2/template"
)

var (
	storesGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		data := dataDefault()

		data["page_title"] = "Stores | " + data["site_title"].(string)

		stores, err := model.StoresAll()
		if err != nil {
			panic(err)
		}

		data["stores"] = stores

		err = template.Views.ExecuteTemplate(w, "stores", "main", data)
		if err != nil {
			panic(err)
		}
	})
)

func init() {
	options.Options.Routes.Get["/stores"] = storesGetHandler
}
