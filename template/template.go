package template

import (
	"github.com/kataras/blocks"
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

var (
	Views *blocks.Blocks
)

const (
	rootDirPath           = "view"
	layoutDirPathRelative = "layout"
)

func init() {
	globalFilepath.Init("..")

	rootDirPathAbs := globalFilepath.Join(rootDirPath)

	Views = blocks.New(rootDirPathAbs).LayoutDir(layoutDirPathRelative)

	err := Views.Load()
	if err != nil {
		panic(err)
	}
}
