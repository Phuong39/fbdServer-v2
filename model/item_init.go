package model

import (
	"encoding/json"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/theTardigrade/fbdServer-v2/database"
)

func init() {
	err := database.BadgerDB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())

			if !strings.HasPrefix(key, itemKeyPrefix) {
				continue
			}

			itemAddToKeyMapAndSlice(key)

			err := item.Value(func(v []byte) error {
				var item2 Item

				err := json.Unmarshal(v, &item2)
				if err != nil {
					return err
				}

				storeAddToItemsMap(item2.StoreName, key)

				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
