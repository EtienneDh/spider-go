# go-crawler
a web crawler in go

**Usage**

./bcrawler -url -domain -depth -csv -private -max

* Url string (default: https://weglot.com)
* Domain string (when not provied, will try to automatically resolve it)
* Depth int (default: 1)
* WriteToCsv bool (default: false)
* Private bool (default: public)
* MaxRequest int (default: -1)
