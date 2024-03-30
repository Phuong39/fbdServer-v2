package model

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/theTardigrade/fbdServer-v2/database"
	hash "github.com/theTardigrade/golang-hash"
)

type Feed struct {
	Url              string    `json:"u"`
	LastDownloadTime time.Time `json:"t"`
}

func FeedKey(url string) []byte {
	var buffer bytes.Buffer

	buffer.WriteString("feed_")
	buffer.WriteString(hash.Uint256String(url).Text(35))

	return buffer.Bytes()
}

func FeedFind(url string) (feed *Feed, found bool, err error) {
	key := FeedKey(url)
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

	var feed2 Feed

	if err = json.Unmarshal(value, &feed2); err != nil {
		return
	}

	feed = &feed2

	return
}

func (f *Feed) Save() (err error) {
	key := FeedKey(f.Url)
	value, err := json.Marshal(f)
	if err != nil {
		return
	}

	err = database.BadgerDB.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value).WithTTL(time.Hour * 24 * 7 * 4) // four weeks
		return txn.SetEntry(entry)
	})
	if err != nil {
		return
	}

	return
}
