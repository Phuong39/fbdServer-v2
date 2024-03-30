package model

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Store string

var (
	storesAllCached       []string
	storesAllCachedInited bool
	storesAllCachedMutex  sync.Mutex
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

	f, err := os.Open(filepath.Join("data", "stores.csv"))
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
