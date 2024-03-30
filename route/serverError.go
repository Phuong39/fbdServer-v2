package route

import (
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/template"
)

var (
	serverErrorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

		err := template.Views.ExecuteTemplate(w, "serverError", "main", nil)
		if err != nil {
			panic(err)
		}
	})
)
