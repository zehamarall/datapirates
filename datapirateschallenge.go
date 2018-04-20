package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//Movie struct
type Movie struct {
	Position string
	Name     string
}

//fetch Load the HTML document
func fetch(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get %s: %v", url, err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {

		if res.StatusCode == http.StatusTooManyRequests {
			return nil, fmt.Errorf("you are being rate limited")
		}

		return nil, fmt.Errorf("bad response from server: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not parse page: %v", err)
	}
	return doc, err
}

//GetListAllGenre Get list all genre in page html
func GetListAllGenre(doc *goquery.Document) (map[string]string, error) {
	listGenre := make(map[string]string)

	doc.Find(".aux-content-widget-2 .subnav_item_main").Each(func(i int, s *goquery.Selection) {
		linkTag := s.Find("a")
		link, _ := linkTag.Attr("href")
		genre := strings.TrimSpace(s.Find("a").Text())
		listGenre[genre] = link
	})
	return listGenre, nil
}

func getGenreRating(page int, genre string, url string) error {

	doc, err := fetch(url)
	if err != nil {
		return err
	}

	fileWriter, err := os.OpenFile("./data/"+genre+".jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	doc.Find(".lister-list .lister-item.mode-advanced .lister-item-content").Each(func(i int, s *goquery.Selection) {
		rating := strings.TrimSpace(s.Find("h3 span").First().Text())
		name := strings.TrimSpace(s.Find("h3 a").Text())
		movie := Movie{
			Position: rating,
			Name:     name,
		}
		//fmt.Printf("# Filme %s - %s  = %s\n", genre, rating, name)
		json.NewEncoder(fileWriter).Encode(movie)
	})

	doc.Find("div.desc a.lister-page-next.next-page").EachWithBreak(func(i int, s *goquery.Selection) bool {
		str, exists := s.Attr("href")
		if exists {
			//Limit 500 films rating
			if page < 9 {
				page++
				urlNext := "http://www.imdb.com/search/title/" + str
				getGenreRating(page, genre, urlNext)
			}
			return false
		}
		return true
	})
	return nil
}

func main() {

	url := "https://www.imdb.com/chart/top"

	// Load the HTML document
	doc, err := fetch(url)
	if err != nil {
		fmt.Printf("could not parse page: %v", err)
	}
	listGenre, err := GetListAllGenre(doc)

	if err != nil {
		log.Fatal(err)
	}
	//var wg sync.WaitGroup

	for genre, link := range listGenre {
		//wg.Add(1)

		go func(genre string, link string) {
			//defer wg.Done()

			err := getGenreRating(0, genre, "http://www.imdb.com"+link)
			if err != nil {
				log.Println(err)
			}
			log.Printf("finished list of %s on jsonl\n", genre)
		}(genre, link)
	}
	err = http.ListenAndServe(":8080", http.FileServer(http.Dir("./data")))

	if err != nil {
		log.Println("could not listen port 8080 ", err)
	}
	log.Println("starting file server on localhost:8080")
	//wg.Wait()
}
