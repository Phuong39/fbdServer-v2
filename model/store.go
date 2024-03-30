package model

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Store string

var (
	StoresAll         []string
	initStoresAllOnce sync.Once
)

func init() {
	initStoresAllOnce.Do(initStoresAll)
}

func initStoresAll() {
	f, err := os.Open(filepath.Join("data", "stores.csv"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, r := range records {
		StoresAll = append(StoresAll, strings.TrimSpace(r[0]))
	}
}
