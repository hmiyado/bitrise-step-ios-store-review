package main

import (
    "bytes"
	"fmt"
	"net/http"
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

func PostToSlack(entry Entry, webhookUrl string) {
	payload := entry.ToSlackPayloadJson()
	http.Post(webhookUrl, "application/json", bytes.NewBuffer([]byte(payload)))
}

func (e *Entry) ToSlackPayloadJson() string {
	rating := ""
	if e.Rating == 0 {
		rating = ":innocent:"
	} else {
		for i := 0; i < e.Rating; i++ {
			rating = rating + ":star: "
		}	
	}
	rating = rating + "("+e.Version+")"

	authorAndDate := fmt.Sprintf("*%s* 、%s", e.Author, e.Updated)

	// https://app.slack.com/block-kit-builder/
	payloadTemplate := `{
		"blocks": [
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "%s",
					"emoji": true
				}
			},
			{
				"type": "section",
				"text": {
					"type": "plain_text",
					"text": "%s",
					"emoji": true
				}
			},
			{
				"type": "section",
				"text": {
					"type": "plain_text",
					"text": "%s",
					"emoji": true
				}
			},
			{
				"type": "context",
				"elements": [
					{
						"type": "mrkdwn",
						"text": "%s"
					}
				]
			}
		]
	}`
	payload := fmt.Sprintf(payloadTemplate, e.Title, rating, e.Content[0].Body, authorAndDate)
	return payload
}
