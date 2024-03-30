package template

import (
	"github.com/kataras/blocks"
)

var (
	Views *blocks.Blocks
)

func init() {
	Views = blocks.New("view").LayoutDir("layout")

	err := Views.Load()
	if err != nil {
		panic(err)
	}
}
