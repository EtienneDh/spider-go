package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"flag"
	"time"
	"net/url"
	"encoding/csv"
	"os"
	"strconv"
)

type Config struct {
	Url string
	Domain string
	Depth int
	WriteToCsv bool
	Private bool
	MaxRequest int
}

// todo dont display already visited link

func main() {	

	Config := getConfig()
	foundUrls := make(map[string]bool)
	requestPerformed := 0

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

		if Config.Private {
			link = link + "?weglot-private=1"
		}

		if Config.MaxRequest != -1 && requestPerformed > Config.MaxRequest {
			return
		}

		if !foundUrls[link] {
			foundUrls[link] = true
			fmt.Println(link)
			e.Request.Visit(link)
			requestPerformed++	
		}
	})

	printInit(Config)
	time.Sleep(2000 * time.Millisecond)

	c.Visit(Config.Url)
	requestPerformed++	
	c.Wait()

	printResults(len(foundUrls))

	if Config.WriteToCsv {
		writeToCsv(foundUrls, Config)
	}
}

func printInit(config Config) {
	fmt.Println("Input url: " + config.Url)
	fmt.Println("Depth: ", config.Depth)	
	fmt.Println("Host", config.Domain)
	fmt.Println("Private mode:", config.Private)
	
	fmt.Println("---------------------------------")
	fmt.Println("Visiting " + config.Url)

}

func printResults(urlsCount int) {
	fmt.Println("---------------------------------")
	fmt.Println("Visited: " + strconv.Itoa(urlsCount) + " urls")
}

func writeToCsv(urls map[string]bool, config Config) {
	// csv init	
	fName := config.Domain + ".csv"
	file, err := os.Create(fName)
	if err != nil {
		panic(err)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	defer fmt.Println("CSV generated at " + dir + "/" + fName)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	i := 0
	for link := range urls {
		writer.Write([]string{strconv.Itoa(i), link})
		i = i + 1
	}
}

func getConfig() Config {
	// Get config from cmd args:
	inputUrl := flag.String("url", "https://weglot.com", "Url to crawl")
	inputDepth := flag.Int("depth", 1, "Depth to crawl")
	allowedDomain := flag.String("domain", "", "Allowed domain")
	writeCsv := flag.Bool("csv", false, "Writes results to CSV")
	private := flag.Bool("private", false, "Crawls with private mode")
	max := flag.Int("max", -1, "Maximum requests to perform")
	flag.Parse()

	config := Config{*inputUrl, *allowedDomain, *inputDepth, *writeCsv, *private, *max}
	config.init()

	return config
}

// Config is passed as pointer to allow modifying the struct
func (config *Config) init() {
	// resolve domain if not passed as arg
	if len(config.Domain) == 0 {
		u, err := url.Parse(config.Url)
	    if err != nil {
	    	fmt.Println("Failed to resolve host")
	        panic(err)
	    }	    
	    host := u.Host
	    config.Domain = host
	}

	// update url for private mode
	if config.Private {
		config.Url = config.Url + "?weglot-private=1"
	}
}

