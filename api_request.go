package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type APIResponse struct {
	Code    int
	Links   []string
	Payload []Word
}

type Word struct {
	T int
	W string
}

const endpoint = "http://node-wordcount-dev.eu-west-3.elasticbeanstalk.com/?url="

func makeWCAPIRequest(url string) (APIResponse, error) {

	var response APIResponse

	if url == "#" || url == "/" {
		return response, errors.New("Cannot crawl this url")
	}

	resp, err := http.Get(endpoint + url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	jsonErr := json.Unmarshal(body, &response)

	if jsonErr != nil {
		log.Println(err)
	}

	return response, nil
}
