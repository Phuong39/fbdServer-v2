package options

import (
	"net/http"

	basicServer "github.com/theTardigrade/golang-basicServer"
)

var (
	Options = basicServer.Options{
		CertFilePath: "certificate/key.cer",
		KeyFilePath:  "certificate/key.key",
		Routes: basicServer.OptionsRoutes{
			Get: map[string]http.Handler{},
		},
	}
)
