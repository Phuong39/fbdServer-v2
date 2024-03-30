package model

import (
	"encoding/json"
	"html/template"
	"math/rand"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/theTardigrade/fbdServer-v2/database"
)

type Item struct {
	StoreName   string        `json:"n"`
	LinkURL     string        `json:"l"`
	ImageURL    string        `json:"i"`
	GUID        string        `json:"g"`
	HashedGUID  string        `json:"h"`
	Title       template.HTML `json:"t"`
	Description template.HTML `json:"d"`
	Keywords    []string      `json:"k"`
	PublishTime time.Time     `json:"p"`
	SetTime     time.Time     `json:"s"`
	Price       uint64        `json:"c"` // US cents
}

func ItemMultipleAtRandom(n int) (items []*Item, err error) {
	itemMap := make(map[*Item]struct{})

	for i := 0; i < (n * 10); i++ {
		if len(itemMap) == n {
			break
		}

		var item *Item

		item, err = ItemAtRandom()
		if err != nil {
			return
		}

		itemMap[item] = struct{}{}
	}

	for item := range itemMap {
		items = append(items, item)
	}

	return
}

func ItemAtRandom() (item *Item, err error) {
	hashedGUID := itemHashedGUIDAtRandom()

	return ItemFromHashedGUID(hashedGUID)
}

func ItemFromHashedGUID(hashedGUID string) (item *Item, err error) {
	var value []byte

	err = database.BadgerDB.View(func(txn *badger.Txn) (err2 error) {
		var rawItem *badger.Item

		rawItem, err2 = txn.Get([]byte(hashedGUID))
		if err2 != nil {
			return err2
		}

		rawItem.Value(func(value2 []byte) error {
			value = value2

			return nil
		})

		return
	})
	if err != nil {
		return
	}

	var item2 Item

	if err = json.Unmarshal(value, &item2); err != nil {
		return
	}

	item = &item2

	return
}

func itemHashedGUIDAtRandom() (guid string) {
	defer itemHashedGUIDMutex.Unlock()
	itemHashedGUIDMutex.Lock()

	index := rand.Intn(len(ItemHashedGUIDSlice))

	return ItemHashedGUIDSlice[index]
}
