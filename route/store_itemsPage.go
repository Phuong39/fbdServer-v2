package route

import (
	htmlTemplate "html/template"
	"math"
	"net/http"
	"strconv"

	"github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
	"github.com/theTardigrade/fbdServer-v2/template"
	basicServer "github.com/theTardigrade/golang-basicServer"
)

const (
	storeItemsPerPage = 48
)

var (
	storeItemsPageGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			}
		}()

		rValues := basicServer.ValuesMapFromRequest(r)

		storeName := rValues["storeName"]
		pageNumberRaw := rValues["pageNumber"]
		pageNumber, err := strconv.Atoi(pageNumberRaw)
		if err != nil || pageNumber < 1 {
			notFoundHandler(w, r)
			return
		}

		data := dataDefault()

		data["store_name"] = storeName
		data["page_title"] = "Page " + pageNumberRaw + " | Store: " + storeName + " | Items | " + data["site_title"].(string)

		items, err := model.ItemMultipleFromStoreName(storeName)
		if err != nil {
			serverErrorHandler(w, r)
			return
		}
		if len(items) == 0 {
			notFoundHandler(w, r)
			return
		}

		pageCount := int(math.Ceil(float64(len(items)) / float64(storeItemsPerPage)))

		data["page_count"] = pageCount
		data["page_number"] = pageNumber
		data["page_description"] = `High-quality customizable items from the ` + storeName + ` store are listed here on ` + data["site_title"].(string) +
			`. Page ` + pageNumberRaw + ` of ` + strconv.Itoa(pageCount) + `.`
		data["page_canonical_link"] = `https://` + data["site_domain"].(string) + `/store/` + storeName + `/items/page/` + pageNumberRaw

		{
			storePageText := htmlTemplate.HTML(`This page contains `)

			if pageCount > 1 {
				storePageText += htmlTemplate.HTML(`part of `)
			}

			storePageText += htmlTemplate.HTML(`our collection ` +
				`of high-quality customizable items from the ` + storeName + ` store.`)

			if pageCount > 1 {
				storePageText += htmlTemplate.HTML("<br />Scroll down and click through to the other pages to see our complete collection.")
			}

			data["store_page_text"] = storePageText
		}

		itemsStartIndex := storeItemsPerPage * (pageNumber - 1)
		itemsEndIndex := storeItemsPerPage * pageNumber

		if itemsStartIndex >= len(items) {
			notFoundHandler(w, r)
			return
		}
		if itemsEndIndex > len(items) {
			itemsEndIndex = len(items)
		}

		prevItemsExist := (itemsStartIndex > 0)
		nextItemsExist := (itemsEndIndex < len(items))

		items = items[itemsStartIndex:itemsEndIndex]

		data["items"] = items

		if prevItemsExist {
			prevPageNumber := pageNumber - 1

			data["prev_page_exists"] = true
			data["prev_page_link"] = `/store/` + storeName + `/items/page/` + strconv.Itoa(prevPageNumber)
		}

		if nextItemsExist {
			nextPageNumber := pageNumber + 1

			data["next_page_exists"] = true
			data["next_page_link"] = `/store/` + storeName + `/items/page/` + strconv.Itoa(nextPageNumber)
		}

		template.Views.ExecuteTemplate(w, "store_page", "main", data)
	})
)

const (
	storeItemsPagePath = "/store/:storeName/items/page/:pageNumber"
)

func init() {
	options.Options.Routes.Get[storeItemsPagePath] = storeItemsPageGetHandler

	stores, err := model.StoresAll()
	if err != nil {
		panic(err)
	}
	for _, s := range stores {
		items, err := model.ItemMultipleFromStoreName(s.Name)
		if err != nil {
			panic(err)
		}
		if len(items) > 0 {
			sitemapPathAdd(`/store/` + s.Name + `/items/page/1`)
		}
		// for i, p := 0, 1; i < len(items); i += storeItemsPerPage {
		// 	sitemapPathAdd(`/store/` + s.Name + `/items/page/` + strconv.Itoa(p))
		// 	p++
		// }
	}
}
