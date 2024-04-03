package route

import (
	"sync"

	"github.com/theTardigrade/fbdServer-v2/environment"
	"github.com/theTardigrade/fbdServer-v2/model"
)

var (
	dataDefaultCached      map[string]interface{}
	dataDefaultCachedMutex sync.RWMutex
)

func dataDefault() (data map[string]interface{}) {
	data = make(map[string]interface{})

	if found := dataDefault_read(data); found {
		return
	}

	dataDefault_readWrite(data)

	return
}

func dataDefault_read(data map[string]interface{}) (found bool) {
	defer dataDefaultCachedMutex.RUnlock()
	dataDefaultCachedMutex.RLock()

	if dataDefaultCached != nil {
		for key, value := range dataDefaultCached {
			data[key] = value
		}

		found = true
		return
	}

	return
}

func dataDefault_readWrite(data map[string]interface{}) {
	defer dataDefaultCachedMutex.Unlock()
	dataDefaultCachedMutex.Lock()

	if dataDefaultCached != nil {
		for key, value := range dataDefaultCached {
			data[key] = value
		}

		return
	}

	environmentKeys := []string{
		"site_domain",
		"site_title",
		// "site_title_initials",
		"referral_site_domain",
		"referral_site_title",
		"referral_query_key",
		"referral_query_value",
	}

	for _, key := range environmentKeys {
		data[key] = environment.Data.MustGet(key)
	}

	{
		stores, err := model.StoresAll()
		if err != nil {
			panic(err)
		}

		storesFiltered := make([]*model.Store, 0, len(stores))

		for _, s := range stores {
			items, err := model.ItemMultipleFromStoreName(s.Name)
			if err != nil {
				panic(err)
			}
			if len(items) > 0 {
				storesFiltered = append(storesFiltered, s)
			}
		}

		data["stores"] = storesFiltered
	}

	{
		dataDefaultCached = make(map[string]interface{})

		for key, value := range data {
			dataDefaultCached[key] = value
		}
	}
}
