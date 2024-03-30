package model

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/dgraph-io/badger/v4"
	"github.com/mmcdole/gofeed"
	"github.com/theTardigrade/fbdServer-v2/database"
	hash "github.com/theTardigrade/golang-hash"
	tasks "github.com/theTardigrade/golang-tasks"
)

const (
	itemTaskInterval = time.Hour * 24 * 7
)

func init() {
	tasks.Set(itemTaskInterval, true, func(id *tasks.Identifier) {
		itemsDownloadAllFromRemoteStore()
	})
}

var (
	itemsDownloadFromRemoteStoreQueryStrings = [...]string{
		"st=popularity&sp=0",
		"st=popularity&sp=1",
		"st=popularity&sp=7",
		"st=popularity&sp=30",
		"st=datecreated",
	}
	itemsDownloadFromRemoteStoreMutex sync.Mutex
)

func itemsDownloadAllFromRemoteStore() (err error) {
	stores := StoresAll[:]

	// shuffle
	sort.Slice(stores, func(i, j int) bool {
		return rand.Intn(2) < 1
	})

	for _, store := range stores {
		queryStrings := itemsDownloadFromRemoteStoreQueryStrings[:]

		// random sort
		sort.Slice(queryStrings, func(i, j int) bool {
			return rand.Float64() >= 0.5
		})

		for _, q := range queryStrings {
			if err = itemsDownloadFromRemoteStore(store, 1, q); err != nil {
				return
			}
		}

		// removeOldItemsFromRemoteStore(store)
	}

	return
}

const (
	itemsDownloadPerFeed = 100
)

func itemsDownloadFromRemoteStore(storeName string, pageNumber int, queryString string) (err error) {
	var totalResults int

	func() (err error) {
		defer itemsDownloadFromRemoteStoreMutex.Unlock()
		itemsDownloadFromRemoteStoreMutex.Lock()

		// foundStore, _, err := StoreLoadOneFromName(storeName)
		// if err != nil || !foundStore {
		// 	panic(err)
		// }

		parser := gofeed.NewParser()
		feedURL := fmt.Sprintf(
			"https://feed.zazzle.com/store/%s/rss?ps=%d&pg=%d&isz=medium",
			storeName,
			itemsDownloadPerFeed,
			pageNumber,
		)
		if queryString != "" {
			feedURL += "&" + queryString
		}

		feed, err := parser.ParseURL(feedURL)
		if err != nil {
			return
		}

		totalResults, err = strconv.Atoi(feed.Extensions["opensearch"]["totalResults"][0].Value)
		if err != nil {
			return
		}

		if len(feed.Items) == 0 {
			return
		}

		// isNewDownloadFeedModel, downloadFeedModel, err := GetDownloadFeedFromURL(feedURL)
		// if err != nil {
		// 	panic(err)
		// } else if !isNewDownloadFeedModel && downloadFeedModel.TooEarly() {
		// 	return
		// }

		for i, l := 0, len(feed.Items); i < l; i++ {
			if err = itemParseFromRemoteStore(feed.Items[i], storeName); err != nil {
				return
			}
		}

		// downloadFeedModel.UpdateLastDownloadTime()
		// downloadFeedModel.Save()

		return
	}()
	if err != nil {
		return
	}

	// if environment.IsDevelopmentMode {
	// 	log.Printf("DOWNLOAD STORE ITEMS: %s (page %d) (%d items per page) (%d results) (%s)\n",
	// 		storeName, pageNumber, itemsPerDownloadFeed, totalResults, queryString)
	// }

	if totalResults > itemsDownloadPerFeed*pageNumber {
		itemsDownloadFromRemoteStore(storeName, pageNumber+1, queryString)
	}

	return
}

func itemParseFromRemoteStore(rawItem *gofeed.Item, storeName string) (err error) {
	linkURL := rawItem.Link
	imageURL := rawItem.Extensions["media"]["content"][0].Attrs["url"]
	guid := rawItem.GUID

	if linkURL == "" || imageURL == "" {
		return
	}

	parsedLinkURL, err := url.Parse(linkURL)
	if err != nil {
		return
	}

	parseLinkURLQuery := parsedLinkURL.Query()

	parseLinkURLQuery.Add("rf", "238606708185114783")

	parsedLinkURL.RawQuery = parseLinkURLQuery.Encode()

	parsedImageURL, err := url.Parse(imageURL)
	if err != nil {
		return
	}

	// foundItem, dbItem, err := LoadItemFromGUID(guid)
	// if err != nil {
	// 	panic(err)
	// }

	publishTime := rawItem.PublishedParsed
	title := template.HTML(rawItem.Title)
	description := template.HTML(rawItem.Extensions["media"]["description"][0].Value)
	keywords := strings.Split(rawItem.Extensions["media"]["keywords"][0].Value, ",")
	priceString := rawItem.Extensions["media"]["price"][0].Value // dollar string

	for i, k := range keywords {
		keywords[i] = strings.TrimSpace(k)
	}

	if publishTime == nil {
		pTime, err := time.Parse(time.RFC1123, rawItem.Published)
		if err != nil {
			return err
		}
		// if time.Since(pTime) <= downloadItemsFromRemoteStoreRecentPublishSkipDuration {
		// 	return // ignore very recent items
		// }
		publishTime = &pTime
	}

	var price uint64
	if priceString != "" {
		var centStringBuilder strings.Builder

		for _, r := range priceString {
			if unicode.IsDigit(r) {
				centStringBuilder.WriteRune(r)
			}
		}

		parsedPrice, err := strconv.ParseUint(centStringBuilder.String(), 10, 64)
		if err == nil {
			price = parsedPrice
		}
	}

	item := Item{
		StoreName:   storeName,
		LinkURL:     parsedLinkURL.String(),
		ImageURL:    parsedImageURL.String(),
		GUID:        guid,
		HashedGUID:  hash.Uint256String(guid).Text(35),
		Title:       title,
		Description: description,
		Keywords:    keywords,
		PublishTime: *publishTime,
		SetTime:     time.Now().UTC(),
		Price:       price,
	}

	itemJSON, err := json.Marshal(item)
	if err != nil {
		return
	}

	err = database.BadgerDB.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(item.HashedGUID), itemJSON).WithTTL(time.Hour * 24 * 7 * 26) // half a year
		return txn.SetEntry(entry)
	})

	itemAddToHashedGUIDMapAndSlice(item.HashedGUID)

	// if foundItem {
	// 	// err = dbItem.Update(
	// 	// 	*parsedLinkURL, *parsedImageURL, storeName, title, description, keywords, price,
	// 	// )
	// } else {
	// 	// dbItem, err = NewItem(
	// 	// 	*parsedLinkURL, *parsedImageURL, guid, storeName, title, description, keywords, *publishTime, price,
	// 	// )
	// }
	// if err != nil {
	// 	panic(err)
	// }

	// if err = dbItem.Save(); err != nil {
	// 	panic(err)
	// }

	// if !foundItem && environment.IsDevelopmentMode {
	// 	log.Printf("ADDED NEW PRODUCT: %s\n", title)
	// }

	return
}
