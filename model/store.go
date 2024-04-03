package model

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

type Store struct {
	Name          string
	FancyName     string
	OnlyUS        bool
	ShowOnSidebar bool
}

var (
	storesAllCached       []*Store
	storesAllCachedInited bool
	storesAllCachedMutex  sync.RWMutex
)

const (
	storesAllFilePath = "data/stores.csv"
)

func StoreFromName(name string) (store *Store, found bool, err error) {
	allStores, err := StoresAll()
	if err != nil {
		return
	}

	for _, store2 := range allStores {
		if store2.Name == name {
			store = store2
			found = true
			return
		}
	}

	return
}

func StoresAll() (stores []*Store, err error) {
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

func storesAll_read(useMutex bool) (stores []*Store, found bool) {
	if useMutex {
		defer storesAllCachedMutex.RUnlock()
		storesAllCachedMutex.RLock()
	}

	if storesAllCachedInited {
		stores = make([]*Store, len(storesAllCached))
		copy(stores, storesAllCached)
		found = true

		return
	}

	return
}

func storesAll_readWrite() (stores []*Store, err error) {
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

	if len(stores) > 0 {
		stores = stores[:0]
	}

	for _, r := range records {
		var onlyUS, showOnSidebar bool

		onlyUS, err = strconv.ParseBool(r[2])
		if err != nil {
			return
		}

		showOnSidebar, err = strconv.ParseBool(r[3])
		if err != nil {
			return
		}

		s := Store{
			Name:          strings.TrimSpace(r[0]),
			FancyName:     strings.TrimSpace(r[1]),
			OnlyUS:        onlyUS,
			ShowOnSidebar: showOnSidebar,
		}

		stores = append(stores, &s)
	}

	sort.Slice(stores, func(i, j int) bool {
		return stores[i].Name < stores[j].Name
	})

	storesAllCached = make([]*Store, len(stores))
	copy(storesAllCached, stores)
	storesAllCachedInited = true

	return
}
