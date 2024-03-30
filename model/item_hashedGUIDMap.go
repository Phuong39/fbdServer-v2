package model

import (
	"sync"

	"github.com/dgraph-io/badger/v4"
	"github.com/theTardigrade/fbdServer-v2/database"
)

var (
	ItemKeyMap              = make(map[string]struct{})
	ItemKeySlice            = make([]string, 0, (1 << 16))
	itemKeyMapAndSliceMutex sync.RWMutex
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
			// var hashedGUID string

			// err := item.Value(func(v []byte) error {
			// 	var item2 Item

			// 	err := json.Unmarshal(v, item2)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	hashedGUID = item2.HashedGUID

			// 	return nil
			// })
			// if err != nil {
			// 	return err
			// }

			itemAddToKeyMapAndSlice(key)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func itemAddToKeyMapAndSlice(key string) {
	defer itemKeyMapAndSliceMutex.Unlock()
	itemKeyMapAndSliceMutex.Lock()

	if _, found := ItemKeyMap[key]; !found {
		ItemKeyMap[key] = struct{}{}
		ItemKeySlice = append(ItemKeySlice, key)
	}
}
