package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type UrlWithCount struct {
	Url       string
	WordCount int
}

func main() {

	config := getConfig()
	foundUrls := make(map[string]bool)
	urlsWithCount := []UrlWithCount{}
	requestPerformed := 0
	// Mutex for requestPerformed
	var mutex = &sync.Mutex{}

	// setup Collector
	c := colly.NewCollector(
		colly.MaxDepth(config.Depth),
		colly.Async(true),
		colly.AllowedDomains(config.Domain),
	)
	c.AllowURLRevisit = false
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// Link Found Callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if config.Private {
			link = link + weglotPrivate
		}

		// Temporary fix to crawl shopify in translated version
		baseUrl := config.Domain
		foundUrl := e.Request.AbsoluteURL(link)
		splitFoundUrl := strings.SplitAfter(foundUrl, baseUrl)

		newUrl := ""
		if len(splitFoundUrl) > 1 {
			newUrl = "/a/l/es" + splitFoundUrl[1]
		} else {
			newUrl = link
		}

		// todo maybe access foundUrls with mutex
		mutex.Lock()
		isAlreadyFound := foundUrls[link]
		mutex.Unlock()
		if link != "" && !isAlreadyFound {
			mutex.Lock()
			foundUrls[link] = true
			mutex.Unlock()
			// crawl for more links
			//e.Request.Visit(e.Request.AbsoluteURL(link))
			e.Request.Visit(newUrl)
		}
	})

	// Request Callback
	c.OnRequest(func(r *colly.Request) {
		// todo : implement crawlRequest & WCRequest, check for both and sum
		if config.MaxRequest != -1 && requestPerformed >= config.MaxRequest {
			return
		}
		mutex.Lock()
		requestPerformed++
		mutex.Unlock()

		url := r.URL.String()
		fmt.Println("Visiting: " + url + " | RP: " + strconv.Itoa(requestPerformed))

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
			urlWithCount := UrlWithCount{url, pageCount}
			urlsWithCount = append(urlsWithCount, urlWithCount)
		}
	})

	printInit(config)
	time.Sleep(2000 * time.Millisecond)

	c.Visit(config.Url)
	c.Wait()

	printResults(urlsWithCount, len(foundUrls), requestPerformed)

	if config.WriteToCsv {
		writeToCsv(urlsWithCount, config.Domain)
	}
}

func printInit(config Config) {
	fmt.Println("hello TEST")
	fmt.Println("Input url: " + config.Url)
	fmt.Println("Depth: ", config.Depth)
	fmt.Println("Host", config.Domain)
	fmt.Println("Private mode:", config.Private)
	fmt.Println("Count words: ", config.Count)

	fmt.Println("---------------------------------")
}

func printResults(urlsWithCount []UrlWithCount, linksFoundCount int, requestPerformed int) {
	totalWordCount := 0
	for _, url := range urlsWithCount {
		totalWordCount = totalWordCount + url.WordCount
	}

	fmt.Println("---------------------------------")
	fmt.Println("Request performed: " + strconv.Itoa(requestPerformed))
	fmt.Println("Number of links found: " + strconv.Itoa(linksFoundCount))
	fmt.Println("Total number of word found: " + strconv.Itoa(totalWordCount))
}
