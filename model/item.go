package model

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/theTardigrade/fbdServer-v2/database"
	"github.com/theTardigrade/fbdServer-v2/random"
	hash "github.com/theTardigrade/golang-hash"
)

type Item struct {
	StoreName   string          `json:"n"`
	LinkURL     string          `json:"l"`
	ImageURL    string          `json:"i"`
	GUID        string          `json:"g"`
	HashedGUID  string          `json:"h"`
	Title       template.HTML   `json:"t"`
	Description template.HTML   `json:"d"`
	Keywords    []template.HTML `json:"k"`
	PublishTime time.Time       `json:"p"`
	SetTime     time.Time       `json:"s"`
	Price       uint64          `json:"c"` // US cents
}

func ItemNew(
	storeName, linkURL, imageURL, guid string,
	title, description template.HTML,
	keywords []template.HTML,
	publishTime time.Time,
	price uint64,
) (item *Item) {
	var hashedGUID string

	{
		hashedGUID = strconv.FormatUint(uint64(hash.Uint32String(guid)), 16)

		for i := len(hashedGUID); i < 4; i++ {
			hashedGUID = "0" + hashedGUID
		}

		hashedGUID = hashedGUID[len(hashedGUID)-4:]

		hashedGUID += strconv.FormatUint(hash.Uint64String(guid), 16)
	}

	item = &Item{
		StoreName:   storeName,
		LinkURL:     linkURL,
		ImageURL:    imageURL,
		GUID:        guid,
		HashedGUID:  hashedGUID,
		Title:       title,
		Description: description,
		Keywords:    keywords,
		PublishTime: publishTime,
		SetTime:     time.Now().UTC(),
		Price:       price,
	}

	return
}

const (
	itemKeyPrefix = "item_"
)

func ItemKey(hashedGUID string) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(itemKeyPrefix)
	buffer.WriteString(hashedGUID)

	return buffer.Bytes()
}

func ItemMultipleAtRandom(n int) (items []*Item, err error) {
	return ItemMultipleAtRandomWithAttempts(n, n*20)
}

func ItemMultipleAtRandomWithAttempts(itemCount, attemptCount int) (items []*Item, err error) {
	itemMap := make(map[*Item]struct{})
	var found bool

	for i := 0; i < attemptCount; i++ {
		if len(itemMap) == itemCount {
			break
		}

		var item *Item

		item, found, err = ItemAtRandom()
		if err != nil {
			return
		}
		if !found {
			continue
		}

		itemMap[item] = struct{}{}
	}

	for item := range itemMap {
		items = append(items, item)
	}

	return
}

func ItemAtRandom() (item *Item, found bool, err error) {
	key, found := itemKeyAtRandom()
	if !found {
		return
	}

	return ItemFromKey([]byte(key))
}

func ItemFromHashedGUID(hashedGUID string) (item *Item, found bool, err error) {
	key := ItemKey(hashedGUID)

	return ItemFromKey(key)
}

func ItemFromKey(key []byte) (item *Item, found bool, err error) {
	var value []byte

	err = database.BadgerDB.View(func(txn *badger.Txn) (err2 error) {
		var rawItem *badger.Item

		rawItem, err2 = txn.Get(key)
		if err2 != nil {
			if err2 == badger.ErrKeyNotFound {
				return nil
			} else {
				return err2
			}
		}

		rawItem.Value(func(value2 []byte) error {
			value = value2
			found = true

			return nil
		})

		return
	})
	if err != nil || !found {
		return
	}

	var item2 Item

	if err = json.Unmarshal(value, &item2); err != nil {
		return
	}

	item = &item2

	return
}

func itemKeyAtRandom() (guid string, found bool) {
	defer itemKeyMapAndSliceMutex.RUnlock()
	itemKeyMapAndSliceMutex.RLock()

	itemsLen := len(ItemKeySlice)

	if (itemsLen) == 0 {
		return
	}

	index := random.Rand.Intn(itemsLen)

	guid = ItemKeySlice[index]
	found = true

	return
}

func ItemMultipleFromStoreName(storeName string) (items []*Item, err error) {
	defer storeItemsMapMutex.RUnlock()
	storeItemsMapMutex.RLock()

	if itemKeysSlice, found := storeItemsMap[storeName]; found {
		for _, key := range *itemKeysSlice {
			var item *Item

			item, found, err = ItemFromKey([]byte(key))
			if err != nil {
				return
			}
			if found {
				items = append(items, item)
			}
		}
	}

	return
}
