package main

import (
	"errors"
	"flag"
	"net/url"
	"strings"
)

// Config is used to parameter the crawler
type Config struct {
	URLS             []string
	Domain           []string
	ExcludedLanguage []string
	Integration      string
	Mode             string
	Private          bool
	Depth            int
	MaxRequest       int
}

var allowedModes = []string{"discover", "crawl"}
var allowedIntegrations = []string{"javascript", "connect", "wordpress-plugin", "shopify-app"}

const defaultURL = ""
const defaultDepth = 5
const defaultPrivate = false
const defaultMaxRequest = 20
const defaultMode = "discover"

const weglotPrivate = "?weglot-private=1"

// NewConfig set up and return a new Config type
func NewConfig() (Config, error) {

	// Get config from cmd args:
	inputURLS := flag.String("url", defaultURL, "Url to crawl")
	inputExcludedLanguages := flag.String("l-excluded", "", "Destination languages")
	integration := flag.String("integration", "", "Project Integration")
	mode := flag.String("mode", defaultMode, "Bot mode: discover or crawl")
	private := flag.Bool("private", defaultPrivate, "Crawls with private mode")
	flag.Parse()

	// extract urls from input
	urls := toArray(*inputURLS)
	if len(urls) == 0 {
		return Config{}, errors.New("You must enter at least 1 url")
	}

	// extract languagesTo
	excludedLanguage := toArray(*inputExcludedLanguages)

	// resolve host
	domain := getDomain(urls[0])
	if domain == "" {
		return Config{}, errors.New("Failed to resolve host")
	}
	// Also add www.domain to go through some redirections
	allDomains := []string{domain, "www." + domain}

	// validate
	if !isValid(*integration, allowedIntegrations) {
		return Config{}, errors.New("You must enter a valid integration: javascript, connect, wordpress plugin or shopify app")
	}
	if !isValid(*mode, allowedModes) {
		return Config{}, errors.New("You must enter a valid mode: discover or crawl")
	}
	// Create struct
	config := Config{
		urls,
		allDomains,
		excludedLanguage,
		*integration,
		*mode,
		*private,
		defaultDepth,
		defaultMaxRequest}

	return config, nil
}

func toArray(s string) []string {

	return strings.Split(s, " ")
}

// Config is passed as pointer to allow modifying the struct
func getDomain(inputURL string) string {
	u, err := url.Parse(inputURL)
	if err != nil {

		return ""
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
