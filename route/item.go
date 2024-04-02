package route

import (
	"html"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-zoo/bone"
	"github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
	"github.com/theTardigrade/fbdServer-v2/template"
	tasks "github.com/theTardigrade/golang-tasks"
)

var (
	itemEscapedKeywordsRegexp = regexp.MustCompile(`(?:&#[0-9]+;|[^a-zA-Z0-9 +-]+)`)
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

		{ // create page_description
			description := html.UnescapeString(string(item.Description))
			description = strings.TrimSpace(description)

			descriptionRunes := []rune(description)
			descriptionRunesLen := len(descriptionRunes)

			for i, r := range descriptionRunes {
				if r == '.' || r == '?' || r == '!' {
					descriptionRunesLen = i + 1
					break
				}

				if i > 0 && (r == '\n' || r == '\r') {
					descriptionRunesLen = i
					break
				}
			}

			description = string(descriptionRunes[:descriptionRunesLen])

			data["page_description"] = description
		}

		itemEscapedKeywords := make([]string, len(item.Keywords))

		for i, keyword := range item.Keywords {
			escapedKeyword := itemEscapedKeywordsRegexp.ReplaceAllString(string(keyword), "")

			itemEscapedKeywords[i] = escapedKeyword
		}

		data["item_escaped_keywords"] = itemEscapedKeywords

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

	func() {
		const itemCount = 10_000

		items, err := model.ItemMultipleAtRandomWithAttempts(itemCount, itemCount)
		if err != nil {
			panic(err)
		}

		paths := make([]string, 0, len(items))

		for _, item := range items {
			path := `/store/` + item.StoreName + `/item/` + item.HashedGUID

			paths = append(paths, path)
		}

		sitemapPathAddMany(paths)
	}()

	tasks.Set(time.Minute*5, false, func(id *tasks.Identifier) {
		if sitemapPathCount() >= 10_000_000 {
			id.Stop()
			return
		}

		const itemCount = 1_000

		items, err := model.ItemMultipleAtRandomWithAttempts(itemCount, itemCount)
		if err != nil {
			panic(err)
		}

		paths := make([]string, 0, len(items))

		for _, item := range items {
			path := `/store/` + item.StoreName + `/item/` + item.HashedGUID

			paths = append(paths, path)
		}

		sitemapPathAddMany(paths)
	})
}
