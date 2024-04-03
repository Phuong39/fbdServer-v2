package model

import "strings"

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
