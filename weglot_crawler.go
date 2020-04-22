package main

// WeglotCrawler is the base struct
type WeglotCrawler struct {
	config Config
}

// NewWeglotCrawler create a new Weglot Crawler
func NewWeglotCrawler() *WeglotCrawler {
	config := getConfig()
	weglotCrawler := WeglotCrawler{config}

	return &weglotCrawler
}

// Run starts crawling & return a collector
func (wCrawler *WeglotCrawler) Run() *Collector {

}
