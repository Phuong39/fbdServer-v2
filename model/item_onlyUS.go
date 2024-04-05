package model

import (
	"encoding/csv"
	"os"
	"sort"
	"strings"

	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

const (
	itemOnlyUSProductSuffixesFilePath = "data/itemOnlyUSSuffixes.csv"
)

var (
	itemOnlyUSProductSuffixes                = []string{} // suffixes must each be at least two bytes long
	itemOnlyUSProductSuffixesLastBytes       = make(map[string]byte)
	itemOnlyUSProductSuffixesSecondLastBytes = make(map[string]byte)
)

func init() {
	err := globalFilepath.Init("..")
	if err != nil {
		panic(err)
	}

	f, err := os.Open(globalFilepath.Join(itemOnlyUSProductSuffixesFilePath))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	for _, line := range lines {
		itemOnlyUSProductSuffixes = append(itemOnlyUSProductSuffixes, line[0])
	}

	for _, s := range itemOnlyUSProductSuffixes {
		if l := len(s); l > 0 {
			itemOnlyUSProductSuffixesLastBytes[s] = s[l-1]

			if l > 1 {
				itemOnlyUSProductSuffixesSecondLastBytes[s] = s[l-2]
			}
		}
	}

	sort.Slice(itemOnlyUSProductSuffixes, func(i, j int) bool {
		sI, sJ := itemOnlyUSProductSuffixes[i], itemOnlyUSProductSuffixes[j]
		bI, bJ := itemOnlyUSProductSuffixesLastBytes[sI], itemOnlyUSProductSuffixesLastBytes[sJ]

		if bI == bJ {
			return itemOnlyUSProductSuffixesSecondLastBytes[sI] < itemOnlyUSProductSuffixesSecondLastBytes[sJ]
		}

		return bI < bJ
	})
}

func (item *Item) OnlyUS() bool {
	lowerTitle := strings.ToLower(string(item.Title))
	lowerTitleLen := len(lowerTitle)
	if lowerTitleLen == 0 {
		return true
	}
	lowerTitleLastByte := lowerTitle[lowerTitleLen-1]
	lowerTitleSecondLastByte := lowerTitle[lowerTitleLen-2]
	lowerTitleShort := lowerTitle[:lowerTitleLen-2]

	var prevLastByte, prevSecondLastByte byte
	for _, s := range itemOnlyUSProductSuffixes {
		lastByte := itemOnlyUSProductSuffixesLastBytes[s]

		if lastByte == prevLastByte {
			continue
		}

		if lastByte == lowerTitleLastByte {
			secondLastByte := itemOnlyUSProductSuffixesSecondLastBytes[s]

			if secondLastByte == prevSecondLastByte {
				continue
			}

			if secondLastByte == lowerTitleSecondLastByte {
				if l := len(s); l > 2 && strings.HasSuffix(lowerTitleShort, s[:l-2]) {
					return true
				}
			} else {
				prevSecondLastByte = secondLastByte
			}
		} else {
			prevLastByte = lastByte
		}
	}

	return false
}
