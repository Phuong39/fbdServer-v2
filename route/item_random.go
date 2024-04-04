package route

import (
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
)

var (
	itemRandomGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		item, found, err := model.ItemAtRandom()
		if err != nil || !found {
			panic(err)
		}

		http.Redirect(w, r, `/store/`+item.StoreName+`/item/`+item.HashedGUID, http.StatusSeeOther)
	})
)

const (
	itemRandomPath = "/item/random"
)

func init() {
	options.Options.Routes.Get[itemRandomPath] = itemRandomGetHandler
}
