package main

import (
	"flag"
	"fmt"
	"web_pull/sources"

	//bwb "web_pull/sources/betterworldbooks"
	gr "web_pull/sources/goodreads"
	//"web_pull/sources/librarything"
)

var SCRAPE_SOURCES = []sources.Source{gr.GoodreadsScraper{}}

func main() {
	isbn := "9780446310789"
	for _, v := range SCRAPE_SOURCES {
		fmt.Println(v.GetReviews(isbn))
	}
}
