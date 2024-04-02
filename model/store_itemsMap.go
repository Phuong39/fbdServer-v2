package model

import "sync"

var (
	storeItemsMap      = make(map[string]*[]string)
	storeItemsMapMutex sync.RWMutex
)

func storeAddToItemsMap(store, itemKey string) {
	defer storeItemsMapMutex.Unlock()
	storeItemsMapMutex.Lock()

	if itemKeysSlice, found := storeItemsMap[store]; !found {
		storeItemsMap[store] = &[]string{itemKey}
	} else {
		shouldAdd := true

		for _, itemKey2 := range *itemKeysSlice {
			if itemKey2 == itemKey {
				shouldAdd = false
				break
			}
		}

		if shouldAdd {
			*itemKeysSlice = append(*itemKeysSlice, itemKey)
		}
	}
}
