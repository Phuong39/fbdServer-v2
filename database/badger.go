package database

import (
	badger "github.com/dgraph-io/badger/v4"
)

var (
	BadgerDB *badger.DB
)

func init() {
	var err error

	BadgerDB, err = badger.Open(badger.DefaultOptions("data/badger"))
	if err != nil {
		panic(err)
	}
}
