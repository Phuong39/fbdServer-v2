package route

import (
	"net/http"

	"github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
)

const (
	itemRandomAttempts = 16
)

var (
	itemRandomGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		var item *model.Item
		var itemFound bool
		var err error

		for i := 0; i < itemRandomAttempts; i++ {
			item, itemFound, err = model.ItemAtRandom()
			if err != nil {
				panic(err)
			}
			if itemFound {
				break
			}
		}
		if !itemFound {
			serverErrorHandler(w, r)
			return
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
