package RSSReader

import (
	RSSReader "RSS-Reader/internal"
	"testing"
)

func TestParse(t *testing.T) {
	var links []string

	result := RSSReader.Parse(links)

	if result != nil {
		t.Errorf("rss parser was passed empty slice of links but did not return nil, in stead returned %v", result)
	}

	/*
		NOTE: This part of the test will crash the test on account of using log.Fatal which uses os.Exit(1)


		links = append(links, "obviously wrong text")
		result = RSSReader.Parse(links)
		if result != nil {
			t.Errorf("rss parser was passed a wrong link, but did not return nil, in stead returned %v", result)
		}
	*/

	links = append(links, "https://rss.com/blog/category/press-releases/feed/")
	links = append(links, "https://blog.jetbrains.com/go/feed")
	result = RSSReader.Parse(links)
	if result == nil {
		t.Error("rss parser was passed a set of correct links, but returned nil")
	}
}
