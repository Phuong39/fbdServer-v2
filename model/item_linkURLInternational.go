package model

import "strings"

func (item *Item) LinkURLInternationalAU() (linkURL string) {
	linkURL = strings.Replace(item.LinkURL, ".com/", ".com.au/", 1) + "&lang=en"

	return
}

func (item *Item) LinkURLInternationalBR() (linkURL string) {
	linkURL = strings.Replace(item.LinkURL, ".com/", ".com.br/", 1) + "&langpt"

	return
}

func (item *Item) LinkURLInternationalCA() (linkURL string) {
	linkURL = strings.Replace(item.LinkURL, ".com/", ".ca/", 1)

	return
}

func (item *Item) LinkURLInternationalDE() (linkURL string) {
	linkURL = strings.Replace(item.LinkURL, ".com/", ".de/", 1) + "&lang=de"

	return
}

func (item *Item) LinkURLInternationalES() (linkURL string) {
	linkURL = strings.Replace(item.LinkURL, ".com/", ".es/", 1) + "&lang=es"

	return
}

func (item *Item) LinkURLInternationalFR() (linkURL string) {
	linkURL = strings.Replace(item.LinkURL, ".com/", ".fr/", 1) + "&lang=fr"

	return
}

func (item *Item) LinkURLInternationalGB() (linkURL string) {
	linkURL = strings.Replace(item.LinkURL, ".com/", ".co.uk/", 1) + "&lang=en"

	return
}

func (item *Item) LinkURLInternationalJP() (linkURL string) {
	linkURL = strings.Replace(item.LinkURL, ".com/", ".co.jp/", 1)

	return
}

func (item *Item) LinkURLInternationalNZ() (linkURL string) {
	linkURL = strings.Replace(item.LinkURL, ".com/", ".co.nz/", 1) + "&lang=en"

	return
}
