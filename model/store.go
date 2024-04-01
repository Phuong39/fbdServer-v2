package model

import (
	"encoding/csv"
	"os"
	"sort"
	"strings"
	"sync"

	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

type Store string

var (
	storesAllCached       []string
	storesAllCachedInited bool
	storesAllCachedMutex  sync.RWMutex
)

const (
	storesAllFilePath = "data/stores.csv"
)

func StoresAll() (stores []string, err error) {
	stores, found := storesAll_read(true)
	if found {
		return
	}

	stores, err = storesAll_readWrite()
	if err != nil {
		return
	}

	return
}

func storesAll_read(useMutex bool) (stores []string, found bool) {
	if useMutex {
		defer storesAllCachedMutex.RUnlock()
		storesAllCachedMutex.RLock()
	}

	if storesAllCachedInited {
		stores = storesAllCached[:]
		return
	}

	return
}

func storesAll_readWrite() (stores []string, err error) {
	defer storesAllCachedMutex.Unlock()
	storesAllCachedMutex.Lock()

	stores, found := storesAll_read(false)
	if found {
		return
	}

	f, err := os.Open(globalFilepath.Join(storesAllFilePath))
	if err != nil {
		return
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return
	}

	for _, r := range records {
		stores = append(stores, strings.TrimSpace(r[0]))
	}

	sort.Strings(stores)

	storesAllCached = stores[:]
	storesAllCachedInited = true

	return
}
