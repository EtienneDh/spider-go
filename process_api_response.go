package main

import "regexp"

func countWords(payload []Word) int {

	// todo Test method (need trim ? removal of blank values ?)
	count := 0
	for _, word := range payload {
		r := regexp.MustCompile("[^\\s]+")
		split := r.FindAllString(word.W, -1)
		count = count + len(split)
	}
	return count
}
