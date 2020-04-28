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
		colly.AllowedDomains(config.Domain...),
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

	url := weglotDiscover.config.URLS[0]
	if weglotDiscover.config.Private {
		url = url + weglotPrivate
	}
	// start crawling and wait until for all thread are done
	weglotDiscover.crawler.Visit(url)
	weglotDiscover.crawler.Wait()

	weglotDiscover.printResults()
}

func (weglotDiscover *WeglotDiscover) printInit() {
	config := weglotDiscover.config
	fmt.Println("Input url: ", config.URLS)
	fmt.Println("Depth: ", config.Depth)
	fmt.Println("Domain", config.Domain)
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

// setCallbacks sets callbacks on link found & request performed events
// @todo Maybe refacto each callback into its own method
func (weglotDiscover *WeglotDiscover) setCallbacks() {
	config := weglotDiscover.config

	// Link Found Callback
	weglotDiscover.crawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// We only want to crawl original version of website when discovering links
		isOriginal, err := IsOriginalLanguage(link, config.Integration, config.ExcludedLanguage)
		if isOriginal && err == nil {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		} else {
			fmt.Println("Not visiting: " + link)
		}

		// need to format url w/ url_formater to handle:
		// adding private part to url
		// crawl shopify url in translated version by adding a/l/_code in url
		if config.Private {
			link = link + weglotPrivate
		}

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
		// fmt.Println("Visited", r.Request.URL)
	})

	weglotDiscover.crawler.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
}
