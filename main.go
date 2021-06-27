package main

import (
    "bytes"
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
	Rating int `xml:"im:rating"`
	Version string `xml:"im:version"`
	Author string `xml:"author>name"`
}

type Content struct {
	XMLNAME xml.Name `xml:"content"`
	Type string `xml:"type,attr"`
	Body string `xml:",innerxml"`
}

func main() {
	appId := os.Getenv("ios_app_id")
	feed := fetchFeed(appId)
	fmt.Printf("count of entries: %d", len(feed.Entries))
	for _, entry := range feed.Entries {
		fmt.Printf(entry.toString())
	}

	webhookUrl := os.Getenv("slack_incoming_webhook_url")
	for i, entry := range feed.Entries {
		postToSlack(entry, webhookUrl)
		if i > 3 {
			break
		}
	}

	// --- Exit codes:
	// The exit code of your Step is very important. If you return
	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// Any non zero exit code will be registered as "failed" by `bitrise`.
	os.Exit(0)
}

func fetchFeed(appId string) Feed {
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

func postToSlack(entry Entry, webhookUrl string) {
	payload := entry.toSlackPayloadJson()
	http.Post(webhookUrl, "application/json", bytes.NewBuffer([]byte(payload)))
}

func (e *Entry) toSlackPayloadJson() string {
	rating := ""
	if e.Rating == 0 {
		rating = ":innocent:"
	} else {
		for i := 0; i < e.Rating; i++ {
			rating += ":star: "
		}	
	}

	// https://app.slack.com/block-kit-builder/
	payloadTemplate := `{
		"blocks": [
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
						"text": "*あやたはわはやかゆ* 、2021/01/04"
					}
				]
			},
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "信頼を裏切られた気分です。",
					"emoji": true
				}
			},
			{
				"type": "section",
				"text": {
					"type": "plain_text",
					"text": "購入してコンビニで支払いを済ませてからの対応が遅すぎませんか？他のネットショッピングサービスでは、例えばAmaz●nでは、購入の手続きが完了してすぐに完了のメールがきます。一方こちらのアプリではメールが届かないだけでなく支払いから一日経っても支払い待ちですと表示されるなど、対応の遅さが非常に目立ちます。金額の大きな買い物だからこそ、この対応の遅さはとても不安になります。アップルの公式だから大丈夫だろうと思って利用したのにとても残念です。ストアだけでなくアップル全体の信用の低下にもつながりかねないと思うので性急に対応すべきです。少なくとも今の状態のままならば私はこのストアをもう一度使いたいとは思えません。改善を期待します。",
					"emoji": true
				}
			}
		]
	}`
	payload := fmt.Sprintf(payloadTemplate, rating)
	return payload
}