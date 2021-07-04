package main

import (
	"os"
	"strconv"
)

func main() {
	appId := os.Getenv("ios_app_id")
	feed := FetchFeed(appId)
	PrintEntries(feed.Entries)

	lastMinutes, _ := strconv.Atoi(os.Getenv("last_minutes"))
	entries := feed.FilterEntriesByLastMinutes(lastMinutes)
	PrintEntries(entries)
	
	webhookUrl := os.Getenv("slack_incoming_webhook_url")
	for _, entry := range entries {
		PostToSlack(entry, webhookUrl)
	}

	// --- Exit codes:
	// The exit code of your Step is very important. If you return
	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// Any non zero exit code will be registered as "failed" by `bitrise`.
	os.Exit(0)
}
