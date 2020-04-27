package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gocolly/colly"
)

// WeglotDiscover is the base struct
type WeglotDiscover struct {
	config           Config
	crawler          *colly.Collector
	mutex            *sync.Mutex
	requestPerformed int
	foundUrls        map[string]bool
	urlsWithCount    []urlWithCount
}

type urlWithCount struct {
	url       string
	wordCount int
}

// NewWeglotDiscover create a new Weglot Crawler that will only extract links from original version of website
func NewWeglotDiscover(config Config) *WeglotDiscover {

	mutex := &sync.Mutex{}
	foundUrls := make(map[string]bool)
	urlsWithCount := []urlWithCount{}

	crawler := colly.NewCollector(
		colly.MaxDepth(config.Depth),
		colly.Async(true),
		colly.AllowedDomains(config.Domain),
	)
	crawler.AllowURLRevisit = false
	crawler.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	weglotDiscover := WeglotDiscover{config, crawler, mutex, 0, foundUrls, urlsWithCount}
	weglotDiscover.setCallbacks()

	return &weglotDiscover
}

// Run starts crawling
func (weglotDiscover *WeglotDiscover) Run() {
	weglotDiscover.printInit()
	weglotDiscover.crawler.Visit(weglotDiscover.config.URLS[0])
	weglotDiscover.crawler.Wait()
	weglotDiscover.printResults()
}

func (weglotDiscover *WeglotDiscover) printInit() {
	config := weglotDiscover.config
	fmt.Println("Input url: ", config.URLS)
	fmt.Println("Depth: ", config.Depth)
	fmt.Println("Host", config.Domain)
	fmt.Println("Private mode:", config.Private)
	fmt.Println("Max requests:", config.MaxRequest)

	fmt.Println("---------------------------------")
}

func (weglotDiscover *WeglotDiscover) printResults() {
	totalWordCount := 0
	for _, url := range weglotDiscover.urlsWithCount {
		totalWordCount = totalWordCount + url.wordCount
	}

	fmt.Println("---------------------------------")
	fmt.Println("Request performed: " + strconv.Itoa(weglotDiscover.requestPerformed))
	fmt.Println("Number of links found: " + strconv.Itoa(len(weglotDiscover.foundUrls)))
	fmt.Println("Total number of word found: " + strconv.Itoa(totalWordCount))
}

// setCallbacks sets callbacks on link find & request performed events
// @todo Maybe refacto each callback into its own method
func (weglotDiscover *WeglotDiscover) setCallbacks() {
	config := weglotDiscover.config

	// Link Found Callback
	weglotDiscover.crawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// need to validate url against languages

		if config.Private {
			link = link + weglotPrivate
		}

		// Todo: see if this logic is really necessary, colly is supposed to handle this by itself
		weglotDiscover.mutex.Lock()
		isAlreadyFound := weglotDiscover.foundUrls[link]
		if link != "" && !isAlreadyFound {
			weglotDiscover.foundUrls[link] = true
			// crawl for more links
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
		weglotDiscover.mutex.Unlock()
	})

	// Request Callback
	weglotDiscover.crawler.OnRequest(func(r *colly.Request) {
		if config.MaxRequest != -1 && weglotDiscover.requestPerformed >= config.MaxRequest {
			return
		}
		weglotDiscover.mutex.Lock()
		weglotDiscover.requestPerformed++
		weglotDiscover.mutex.Unlock()

		url := r.URL.String()
		fmt.Println("Visiting: " + url + " | RP: " + strconv.Itoa(weglotDiscover.requestPerformed))
	})

	weglotDiscover.crawler.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})
}
