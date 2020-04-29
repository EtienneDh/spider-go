package main

import (
	"fmt"
	"os"
)

func main() {
	// Maybe some type of configs will have an effect here
	// Like if we're in generate content mode and receive a list of URLs
	// It may be easier to run len(URLS) times the bot

	// To achieve this, I would need to refacto the config out of weglotCrawler.initialize
	// and instead pass it as args to New method

	config, err := NewConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if config.Mode == "discover" {
		weglotDiscover := NewWeglotDiscover(config)
		weglotDiscover.Run()
	}

	// weglotCrawler := NewWeglotCrawler()
	// weglotCrawler.Run()
}
