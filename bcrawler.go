package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"flag"
	"time"
	"net/url"
)

type Url struct {
	Name string
}

func main() {	

	inputUrl := flag.String("url", "https://weglot.com", "Url to crawl")
	inputDepth := flag.Int("depth", 1, "Depth to crawl")
	allowedHost := flag.String("host", "", "Allowed host")

	flag.Parse()

	if(len(*allowedHost) == 0) {
		u, err := url.Parse(*inputUrl)
	    if err != nil {
	        panic(err)
	    }
	    fmt.Println(u.Host)
	    host := u.Host
	    allowedHost = &host
	}

	fmt.Println("Input url: " + *inputUrl)
	fmt.Println("Depth: ", *inputDepth)	
	fmt.Println("Host", *allowedHost)

	time.Sleep(2000 * time.Millisecond)

	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 2, so only the links on the scraped page
		// and links on those pages are visited
		colly.MaxDepth(*inputDepth),
		colly.Async(true),		
		colly.AllowedDomains(*allowedHost),

	)
	c.AllowURLRevisit = false

	fmt.Println("---------------------------------")
	fmt.Println("Starting...")
	fmt.Println("---------------------------------")

	urls := make([]Url, 0, 200)

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		url := Url{link}	

		// Print link
		fmt.Println(link)
		// Visit link found on page on a new thread
		e.Request.Visit(link)

		urls = append(urls, url)
	})

	fmt.Println("Visiting " + *inputUrl)
	c.Visit(*inputUrl)	
	c.Wait()

	fmt.Println("---------------------------------")
	fmt.Println("Visisted:")
	fmt.Println(len(urls))
	fmt.Println("urls")
}

