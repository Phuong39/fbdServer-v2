package route

import (
	"github.com/theTardigrade/fbdServer-v2/environment"
)

func dataDefault() (data map[string]interface{}) {
	data = make(map[string]interface{})

	environmentKeys := []string{
		"site_title",
		"referral_site_domain",
		"referral_site_title",
		"referral_query_key",
		"referral_query_value",
	}

	for _, key := range environmentKeys {
		data[key] = environment.Data.MustGet(key)
	}

	return
}
