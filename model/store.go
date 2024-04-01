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
	storesAllCachedMutex  sync.Mutex
)

const (
	storesAllFilePath = "data/stores.csv"
)

func StoresAll() (stores []string, err error) {
	if storesAllCachedInited {
		stores = storesAllCached[:]
		return
	}

	defer storesAllCachedMutex.Unlock()
	storesAllCachedMutex.Lock()

	if storesAllCachedInited {
		stores = storesAllCached[:]
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
