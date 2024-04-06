package model

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/theTardigrade/fbdServer-v2/database"
)

func (f *Feed) Save() (err error) {
	value, err := json.Marshal(f)
	if err != nil {
		return
	}

	key := FeedKey(f.Url)

	err = database.BadgerDB.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value).WithTTL(time.Hour * 24 * 7 * 4) // four weeks
		return txn.SetEntry(entry)
	})
	if err != nil {
		return
	}

	return
}
