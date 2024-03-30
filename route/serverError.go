package route

import (
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/template"
)

var (
	serverErrorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"message": "The server has experienced an unexpected error.",
		}

		w.WriteHeader(http.StatusInternalServerError)

		err := template.Views.ExecuteTemplate(w, "serverError", "main", data)
		if err != nil {
			panic(err)
		}
	})
)
