package model

import (
	"sync"

	"github.com/dgraph-io/badger/v4"
	"github.com/theTardigrade/fbdServer-v2/database"
)

var (
	ItemHashedGUIDMap   = make(map[string]struct{})
	ItemHashedGUIDSlice = make([]string, 0, (1 << 16))
	itemHashedGUIDMutex sync.Mutex
)

func init() {
	err := database.BadgerDB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			hashedGUID := string(item.Key())
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

			itemAddToHashedGUIDMapAndSlice(hashedGUID)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func itemAddToHashedGUIDMapAndSlice(hashedGUID string) {
	defer itemHashedGUIDMutex.Unlock()
	itemHashedGUIDMutex.Lock()

	if _, found := ItemHashedGUIDMap[hashedGUID]; !found {
		ItemHashedGUIDMap[hashedGUID] = struct{}{}
		ItemHashedGUIDSlice = append(ItemHashedGUIDSlice, hashedGUID)
	}
}
