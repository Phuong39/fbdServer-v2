package model

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Store string

var (
	storesAllCached       []string
	storesAllCachedInited bool
	storesAllCachedMutex  sync.Mutex
)

func StoresAll() (stores []string) {
	fmt.Println(storesAllCachedInited, storesAllCached)
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
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, r := range records {
		stores = append(stores, strings.TrimSpace(r[0]))
	}

	storesAllCached = stores[:]
	storesAllCachedInited = true

	return
}
