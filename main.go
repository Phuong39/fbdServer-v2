package main

import (
	"github.com/theTardigrade/fbdServer-v2/database"
	_ "github.com/theTardigrade/fbdServer-v2/middleware"
	_ "github.com/theTardigrade/fbdServer-v2/minify"
	_ "github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
	_ "github.com/theTardigrade/fbdServer-v2/route"
	basicServer "github.com/theTardigrade/golang-basicServer"
)

func main() {
	defer database.Close()

	if err := basicServer.Serve(options.Options); err != nil {
		panic(err)
	}
}
