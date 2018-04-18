package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestGetListAllGenre(t *testing.T) {
	dat, _ := ioutil.ReadFile("data/htmlCodeListallgenre.html")

	html := string(dat)

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader((html)))

	m, err := GetListAllGenre(doc)

	if err != nil {
		t.Error("var error is nil ", err)
	}

	if len(m) != 21 {
		t.Error("Number of  genre was incorrect ", len(m))
	}

}
