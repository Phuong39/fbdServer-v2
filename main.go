package main

import (
	"log"

	"github.com/theTardigrade/fbdServer-v2/database"
	_ "github.com/theTardigrade/fbdServer-v2/middleware"
	_ "github.com/theTardigrade/fbdServer-v2/model"
	"github.com/theTardigrade/fbdServer-v2/options"
	_ "github.com/theTardigrade/fbdServer-v2/route"
	basicServer "github.com/theTardigrade/golang-basicServer"
)

func main() {
	defer database.Close()

	basicServer.ServeContinuously(options.Options, func(err error) {
		log.Println(err)
	})
}
