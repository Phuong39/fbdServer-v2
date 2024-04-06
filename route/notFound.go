package route

import (
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/options"
	"github.com/theTardigrade/fbdServer-v2/template"
)

var (
	notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		data := dataDefault()

		data["heading"] = "Error 404"
		data["page_title"] = data["heading"].(string) + " | " + data["site_title"].(string)
		data["message"] = "The page you are looking for cannot be found."

		w.WriteHeader(http.StatusNotFound)

		template.Views.ExecuteTemplate(w, "notFound", "main", data)
	})
)

func init() {
	options.Options.Routes.NotFound = notFoundHandler
}
