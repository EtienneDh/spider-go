package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"
)

type Config struct {
	URLS        []string
	Domain      string
	Private     bool
	Integration string
	Mode        string
}

var allowedModes = []string{"discover", "crawl"}
var allowedIntegrations = []string{"javascript", "connect", "wordpress plugin", "shopify app"}

const defaultURL = "https://weglot.com"
const defaultDepth = 5
const defaultDomain = ""
const defaultPrivate = false
const defaultMaxRequest = -1
const defaultCount = false

const weglotPrivate = "?weglot-private=1"

// NewConfig set up and return a new Config type
func NewConfig() (Config, error) {

	// Get config from cmd args:
	inputURLS := flag.String("url", defaultURL, "Url to crawl")
	private := flag.Bool("private", defaultPrivate, "Crawls with private mode")
	integration := flag.String("integration", "", "Project Integration")
	mode := flag.String("mode", "discover", "Bot mode: discover or crawl")
	flag.Parse()

	// extract urls
	urls := getURLS(*inputURLS)
	if len(urls) == 0 {
		return Config{}, errors.New("You must enter at least 1 url")
	}
	// resolve host
	domain := getDomain(urls[0])

	// validate
	if !isValid(*integration, allowedIntegrations) {
		return Config{}, errors.New("You must enter a valid integration; javascript, connect, wordpress plugin or shopify app")
	}
	if !isValid(*mode, allowedModes) {
		return Config{}, errors.New("You must enter a valid mode: discover or crawl")
	}
	config := Config{urls, domain, *private, *integration, *mode}

	return config, nil
}

func getURLS(urlsAsString string) []string {

	return strings.Split(urlsAsString, " ")
}

// Config is passed as pointer to allow modifying the struct
func getDomain(inputURL string) string {
	u, err := url.Parse(inputURL)
	if err != nil {
		fmt.Println("Failed to resolve host")
		panic(err)
	}

	return u.Host
}

func isValid(toValidate string, authorizedValues []string) bool {
	for _, value := range authorizedValues {
		if toValidate == value {
			return true
		}
	}

	return false
}
