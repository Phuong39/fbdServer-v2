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

		w.WriteHeader(http.StatusNotFound)

		err := template.Views.ExecuteTemplate(w, "notFound", "main", nil)
		if err != nil {
			panic(err)
		}
	})
)

func init() {
	options.Options.Routes.NotFound = notFoundHandler
}
