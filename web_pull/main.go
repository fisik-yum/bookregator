package main

import (
	gr "web_pull/sources/goodreads"
	hc"web_pull/sources/hardcover"
	"web_pull/sources"
)

var SCRAPE_SOURCES=[2]sources.Source{gr.GoodreadsScraper{},hc.HardcoverScraper{}}

func main(){
	//h,_:=gr.GoodreadsScraper{}.GetReviews("9781439550410")
	//fmt.Println(h[0].Rating)
}
