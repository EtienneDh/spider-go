package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func writeToCsv(urlsWithCount []UrlWithCount, domain string) {
	// csv init
	fName := domain + ".csv"
	file, err := os.Create(fName)
	if err != nil {
		panic(err)
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
	for _, urlWithCount := range urlsWithCount {
		writer.Write([]string{urlWithCount.Url, strconv.Itoa(urlWithCount.WordCount)})
		i = i + 1
	}
}
