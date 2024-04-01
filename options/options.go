package options

import (
	"net/http"

	basicServer "github.com/theTardigrade/golang-basicServer"
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
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

func init() {
	globalFilepath.Init("..")

	Options.CertFilePath = globalFilepath.Join(Options.CertFilePath)
	Options.KeyFilePath = globalFilepath.Join(Options.KeyFilePath)
}
