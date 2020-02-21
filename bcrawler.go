package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"flag"
	"time"
	"net/url"
	"encoding/csv"
	"os"
)


type Config struct {
	Url string
	Domain string
	Depth int
	WriteToCsv bool
}

// todo dont display already visited link

func main() {	

	Config := getConfig()
	foundUrls := make(map[string]bool)

	// setup Collector
	c := colly.NewCollector(		
		colly.MaxDepth(Config.Depth),
		colly.Async(true),		
		colly.AllowedDomains(Config.Domain),

	)
	c.AllowURLRevisit = false	
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if !foundUrls[link] {
			foundUrls[link] = true
			fmt.Println(link)
			e.Request.Visit(link)	
		} 
		
	})

	printInit(Config)
	time.Sleep(2000 * time.Millisecond)

	c.Visit(Config.Url)	
	c.Wait()

	printResults(foundUrls)

	if Config.WriteToCsv {
		writeToCsv(foundUrls, Config)
	}

}

func printInit(config Config) {
	fmt.Println("Input url: " + config.Url)
	fmt.Println("Depth: ", config.Depth)	
	fmt.Println("Host", config.Domain)
	
	fmt.Println("---------------------------------")
	fmt.Println("Visiting " + config.Url)

}

func printResults(urls map[string]bool) {
	fmt.Println("---------------------------------")
	fmt.Println("Visisted:")
	fmt.Println(len(urls))
	fmt.Println("urls")
}

func writeToCsv(urls map[string]bool, config Config) {
	// csv init	
	fName := config.Domain + ".csv"
	file, err := os.Create(fName)
	if err != nil {
		panic(err)
		return
	}

	fmt.Println("Writing to " + config.Domain + ".csv ...")

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	defer fmt.Println("done")

	for link := range urls {
		writer.Write([]string{link})
	}
}

func getConfig() Config {
	// Get config from cmd args:
	inputUrl := flag.String("url", "https://weglot.com", "Url to crawl")
	inputDepth := flag.Int("depth", 1, "Depth to crawl")
	allowedDomain := flag.String("domain", "", "Allowed domain")
	writeCsv := flag.Bool("csv", false, "Writes results to CSV")
	flag.Parse()

	// Resolve domain if not passed as arg
	if(len(*allowedDomain) == 0) {
		u, err := url.Parse(*inputUrl)
	    if err != nil {
	        panic(err)
	    }
	    fmt.Println(u.Host)
	    host := u.Host
	    allowedDomain = &host
	}

	config := Config{*inputUrl, *allowedDomain, *inputDepth, *writeCsv}

	return config
}

