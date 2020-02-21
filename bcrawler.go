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

	// Get config from cmd args:
	inputUrl := flag.String("url", "https://weglot.com", "Url to crawl")
	inputDepth := flag.Int("depth", 1, "Depth to crawl")
	allowedDomain := flag.String("domain", "", "Allowed domain")

	flag.Parse()

	// 
	if(len(*allowedDomain) == 0) {
		u, err := url.Parse(*inputUrl)
	    if err != nil {
	        panic(err)
	    }
	    fmt.Println(u.Host)
	    host := u.Host
	    allowedDomain = &host
	}

	fmt.Println("Input url: " + *inputUrl)
	fmt.Println("Depth: ", *inputDepth)	
	fmt.Println("Host", *allowedDomain)

	time.Sleep(2000 * time.Millisecond)

	// setup Collector
	c := colly.NewCollector(		
		colly.MaxDepth(*inputDepth),
		colly.Async(true),		
		colly.AllowedDomains(*allowedDomain),

	)
	c.AllowURLRevisit = false

	urls := make([]Url, 0, 200)

	fmt.Println("---------------------------------")
	fmt.Println("Starting...")
	fmt.Println("---------------------------------")	
	
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		url := Url{link}
		
		fmt.Println(link)
		
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

