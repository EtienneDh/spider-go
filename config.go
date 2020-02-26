package main

import (
	"flag"
	"fmt"
	"net/url"
)

type Config struct {
	Url        string
	Domain     string
	Depth      int
	WriteToCsv bool
	Private    bool
	MaxRequest int
}

const defaultURL = "https://weglot.com"
const defaultDepth = 1
const defaultDomain = ""
const defaultCSV = false
const defaultPrivate = false
const defaultMaxRequest = -1

const weglotPrivate = "?weglot-private=1"

func getConfig() Config {
	// Get config from cmd args:
	inputUrl := flag.String("url", defaultURL, "Url to crawl")
	inputDepth := flag.Int("depth", defaultDepth, "Depth to crawl")
	allowedDomain := flag.String("domain", defaultDomain, "Allowed domain")
	writeCsv := flag.Bool("csv", defaultCSV, "Writes results to CSV")
	private := flag.Bool("private", defaultPrivate, "Crawls with private mode")
	max := flag.Int("max", defaultMaxRequest, "Maximum requests to perform")
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
		config.Url = config.Url + weglotPrivate
	}
}
