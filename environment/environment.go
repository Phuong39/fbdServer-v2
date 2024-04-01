package environment

import (
	envStore "github.com/theTardigrade/golang-envStore"
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

const (
	filePath = "data/environment.env"
)

var (
	Data   *envStore.Environment
	config = &envStore.Config{
		FromFilePaths: []string{
			globalFilepath.Join(filePath),
		},
		IgnoreEmptyLines: true,
		UseMutex:         true,
	}
)

var (
	IsDevelopmentMode bool
	IsProductionMode  bool
)

func init() {
	var err error

	Data, err = envStore.New(config)
	if err != nil {
		panic(err)
	}

	// in production mode, unless otherwise stated
	IsDevelopmentMode = Data.LazyGetBool("enable_development_mode")
	IsProductionMode = !IsDevelopmentMode
}

func IsKeyNotFoundErr(err error) bool {
	return err == envStore.ErrKeyNotFound
}
