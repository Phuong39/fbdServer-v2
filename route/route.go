package route

import (
	"sync"

	"github.com/theTardigrade/fbdServer-v2/environment"
)

var (
	dataDefaultCached      map[string]interface{}
	dataDefaultCachedMutex sync.Mutex
)

func dataDefault() (data map[string]interface{}) {
	data = make(map[string]interface{})

	if dataDefaultCached != nil {
		for key, value := range dataDefaultCached {
			data[key] = value
		}

		return
	}

	environmentKeys := []string{
		"site_domain",
		"site_title",
		"referral_site_domain",
		"referral_site_title",
		"referral_query_key",
		"referral_query_value",
	}

	for _, key := range environmentKeys {
		data[key] = environment.Data.MustGet(key)
	}

	defer dataDefaultCachedMutex.Unlock()
	dataDefaultCachedMutex.Lock()

	dataDefaultCached = data

	return
}
