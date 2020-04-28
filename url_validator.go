package main

import (
	"fmt"
	"regexp"
)

// IsOriginalLanguage returns trues if url does not have any language code
func IsOriginalLanguage(url string, integration string, languagesTo []string) bool {
	pattern := ""

	switch integration {
	case "wordpress-plugin":
		pattern = getRegexPatternForWordPress(languagesTo)
	case "shopify-app":
		pattern = getRegexPatternForJavascript(languagesTo)
	case "javascript":
		pattern = getRegexPatternForJavascript(languagesTo)
	default:
		pattern = getRegexPatternForWordPress(languagesTo)
	}

	match, err := regexp.Match(pattern, []byte(url))
	if err != nil {
		fmt.Println(err)
	}

	return !match
}

// Looks for mywebsite/_code/resource in URL or mywebsite/resource/_code
// Generates regex like: /(en|es|ja)/|/(en|es|ja)$
func getRegexPatternForWordPress(languagesTo []string) string {
	pattern := "/("
	pattern = addLanguageCodes(pattern, languagesTo)
	pattern += ")/|/("
	pattern = addLanguageCodes(pattern, languagesTo)
	pattern += ")$"

	return pattern
}

// Looks for mywebsite/a/l/_code/resource
// Generates regex like: /a/l/(en|es|ja)/
func getRegexPatternForJavascript(languagesTo []string) string {
	pattern := "/a/l/("
	pattern = addLanguageCodes(pattern, languagesTo)
	pattern += ")/"

	return pattern
}

func addLanguageCodes(pattern string, languagesTo []string) string {
	for k, v := range languagesTo {
		pattern += v
		if k != len(languagesTo)-1 {
			pattern += "|"
		}
	}

	return pattern
}
