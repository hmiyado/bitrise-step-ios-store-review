package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	appId := os.Getenv("ios_app_id")
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
 
	fmt.Println(string(body))

	//
	// --- Step Outputs: Export Environment Variables for other Steps:
	// You can export Environment Variables for other Steps with
	//  envman, which is automatically installed by `bitrise setup`.
	// A very simple example:
	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "EXAMPLE_STEP_OUTPUT", "--value", "the value you want to share").CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}
	// You can find more usage examples on envman's GitHub page
	//  at: https://github.com/bitrise-io/envman

	//
	// --- Exit codes:
	// The exit code of your Step is very important. If you return
	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// Any non zero exit code will be registered as "failed" by `bitrise`.
	os.Exit(0)
}
