package model

import (
	"bytes"
	"encoding/json"
	"html/template"
	"math/rand"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/theTardigrade/fbdServer-v2/database"
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

func ItemKey(hashedGUID string) []byte {
	var buffer bytes.Buffer

	buffer.WriteString("item_")
	buffer.WriteString(hashedGUID)

	return buffer.Bytes()
}

func ItemMultipleAtRandom(n int) (items []*Item, err error) {
	itemMap := make(map[*Item]struct{})
	var found bool

	for i := 0; i < (n * 10); i++ {
		if len(itemMap) == n {
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
	key := itemKeyAtRandom()

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

func itemKeyAtRandom() (guid string) {
	defer itemKeyMapAndSliceMutex.RUnlock()
	itemKeyMapAndSliceMutex.RLock()

	index := rand.Intn(len(ItemKeySlice))

	return ItemKeySlice[index]
}
