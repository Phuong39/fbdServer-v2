package main

import (
	"log"
	"os"

	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

func init() {
	globalFilepath.Init()

	log.SetOutput(os.Stderr)
}
