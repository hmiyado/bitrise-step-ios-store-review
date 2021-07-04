package main

import (
    "bytes"
	"fmt"
	"net/http"
)

func PostToSlack(entry Entry, webhookUrl string) {
	payload := entry.toSlackPayloadJson()
	http.Post(webhookUrl, "application/json", bytes.NewBuffer([]byte(payload)))
}

func (e *Entry) toSlackPayloadJson() string {
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
