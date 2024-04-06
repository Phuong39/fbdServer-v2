package model

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v4"
	"github.com/theTardigrade/fbdServer-v2/database"
)

func (item *Item) Save() (err error) {
	value, err := json.Marshal(item)
	if err != nil {
		return
	}

	key := ItemKey(item.HashedGUID)

	err = database.BadgerDB.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value) //.WithTTL(time.Hour * 24 * 7 * 26) // half a year
		return txn.SetEntry(entry)
	})
	if err != nil {
		return
	}

	keyString := string(key)

	itemAddToKeyMapAndSlice(keyString)
	storeAddToItemsMap(item.StoreName, keyString)

	return
}
