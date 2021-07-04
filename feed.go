package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Feed struct {
	XMLNAME xml.Name `xml:"feed"`
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	XMLNAME xml.Name `xml:"entry"`
	Updated time.Time `xml:"updated"`
	Title string `xml:"title"`
	Content []Content `xml:"content"` 
	Rating int `xml:"rating"`
	Version string `xml:"version"`
	Author string `xml:"author>name"`
}

type Content struct {
	XMLNAME xml.Name `xml:"content"`
	Type string `xml:"type,attr"`
	Body string `xml:",innerxml"`
}

func FetchFeed(appId string) Feed {
	appUrl := "https://itunes.apple.com/jp/rss/customerreviews/page=1/id=" + appId + "/sortby=mostrecent/xml?urlDesc=/customerreviews/id=" + appId + "/sortBy=mostRecent/json"

	response, err := http.Get(appUrl)
	defer response.Body.Close()
	if err != nil {
		fmt.Printf("Failed to get review from %s, error: %#v\n", appUrl, err)
		os.Exit(1)
	} 
	responseStatus := string(response.Status)
	if !strings.Contains(responseStatus, "200") {
		fmt.Printf("Response from ios store is not OK but: %s", responseStatus)
		os.Exit(1)
	}
 
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body from %s, error: %#v \n", appUrl, err)
		os.Exit(1)
	}

	feed := Feed{}
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		fmt.Printf("Fialed to parse xml of store review,\n%s\n", string(body))
		os.Exit(1)
	}
	return feed
}

func (e *Entry) toString() string {
	return fmt.Sprintf("[%s]<V:%s><R:%d> %s -- %s\n%s\n", e.Updated, e.Version, e.Rating, e.Title, e.Author, e.Content[0].Body)
}

func PrintEntries(entries []Entry) {
	fmt.Printf("count of entries: %d\n", len(entries))
	for _, entry := range entries {
		fmt.Printf(entry.toString())
	}
}

func (f Feed) FilterEntriesByLastMinutes(lastMinutes int) []Entry {
	entries := []Entry{}
	lastValidTime := time.Now().Truncate(time.Duration(lastMinutes) * time.Minute)

	for _, entry := range f.Entries {
		if entry.Updated.After(lastValidTime) {
			entries = append(entries, entry)
		}
	}
	return entries
}
