package database

import (
	badger "github.com/dgraph-io/badger/v4"
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

const (
	badgerDataPath = "data/badger"
)

var (
	BadgerDB *badger.DB
)

func init() {
	globalFilepath.Init("..")

	badgerDataPathAbs := globalFilepath.Join(badgerDataPath)

	var err error

	BadgerDB, err = badger.Open(badger.DefaultOptions(badgerDataPathAbs))
	if err != nil {
		panic(err)
	}
}
