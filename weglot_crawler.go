package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gocolly/colly"
)

// WeglotCrawler is the base struct
type WeglotCrawler struct {
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

// NewWeglotCrawler create a new Weglot Crawler
func NewWeglotCrawler() *WeglotCrawler {
	config := getConfig()
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

	weglotCrawler := WeglotCrawler{config, crawler, mutex, 0, foundUrls, urlsWithCount}
	weglotCrawler.setCallbacks()

	return &weglotCrawler
}

// Run starts crawling & return a collector
func (weglotCrawler *WeglotCrawler) Run() {
	weglotCrawler.printInit()
	weglotCrawler.crawler.Visit(weglotCrawler.config.Url)
	weglotCrawler.crawler.Wait()
	weglotCrawler.printResults()
}

func (weglotCrawler *WeglotCrawler) printInit() {
	config := weglotCrawler.config
	fmt.Println("Input url: " + config.Url)
	fmt.Println("Depth: ", config.Depth)
	fmt.Println("Host", config.Domain)
	fmt.Println("Private mode:", config.Private)
	fmt.Println("Count words: ", config.Count)

	fmt.Println("---------------------------------")
}

func (weglotCrawler *WeglotCrawler) printResults() {
	totalWordCount := 0
	for _, url := range weglotCrawler.urlsWithCount {
		totalWordCount = totalWordCount + url.wordCount
	}

	fmt.Println("---------------------------------")
	fmt.Println("Request performed: " + strconv.Itoa(weglotCrawler.requestPerformed))
	fmt.Println("Number of links found: " + strconv.Itoa(len(weglotCrawler.foundUrls)))
	fmt.Println("Total number of word found: " + strconv.Itoa(totalWordCount))
}

// setCallbacks sets callbacks on link find & request performed events
// @todo Maybe refacto each callback into its own method
func (weglotCrawler *WeglotCrawler) setCallbacks() {
	config := weglotCrawler.config

	// Link Found Callback
	weglotCrawler.crawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if config.Private {
			link = link + weglotPrivate
		}

		weglotCrawler.mutex.Lock()
		isAlreadyFound := weglotCrawler.foundUrls[link]
		if link != "" && !isAlreadyFound {
			weglotCrawler.foundUrls[link] = true
			// crawl for more links
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
		weglotCrawler.mutex.Unlock()
	})

	// Request Callback
	weglotCrawler.crawler.OnRequest(func(r *colly.Request) {
		if config.MaxRequest != -1 && weglotCrawler.requestPerformed >= config.MaxRequest {
			return
		}
		weglotCrawler.mutex.Lock()
		weglotCrawler.requestPerformed++
		weglotCrawler.mutex.Unlock()

		url := r.URL.String()
		fmt.Println("Visiting: " + url + " | RP: " + strconv.Itoa(weglotCrawler.requestPerformed))

		// get wordcount for this url
		if config.Count {
			apiResponse, error := makeWCAPIRequest(url)
			pageCount := 0
			if error == nil {
				pageCount = countWords(apiResponse.Payload)
				fmt.Println("---------------------------------")
				fmt.Println(url + " has " + strconv.Itoa(pageCount) + " words")
				fmt.Println("---------------------------------")
			}
			urlWithCount := urlWithCount{url, pageCount}
			weglotCrawler.urlsWithCount = append(weglotCrawler.urlsWithCount, urlWithCount)
		}
	})
}
