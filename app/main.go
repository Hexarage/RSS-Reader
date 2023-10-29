package main

import (
	RSSReader "RSS-Reader/internal"
	"fmt"
)

func main() {
	fmt.Println("Hello there")
	test := RSSReader.RSSItem{}
	test.Parse(nil)
}
