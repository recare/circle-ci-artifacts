package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	user := os.Args[1]
	project := os.Args[2]
	buildNumber := os.Args[3]
	binaryName := os.Args[4]
	destination := os.Args[5]
	token := os.Args[6]

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://circleci.com/api/v1.1/project/github/%s/%s/%s/artifacts?circle-token=%s", user, project, buildNumber, token), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var artifacts []Artifact
	if err := json.Unmarshal(body, &artifacts); err != nil {
		panic(err)
	}

	for _, artifact := range artifacts {
		if artifact.PrettyPath == fmt.Sprintf("$CIRCLE_ARTIFACTS/%s", binaryName) {
			binary, err := http.Get(fmt.Sprintf("%s?circle-token=%s", artifact.Url, token))

			if err != nil {
				panic(err)
			}

			defer binary.Body.Close()
			out, err := os.Create(fmt.Sprintf("%s/%s", destination, binaryName))
			if err != nil {
				panic(err)
			}

			defer out.Close()
			io.Copy(out, resp.Body)
		}

		return
	}
}

type Artifact struct {
	Path       string `json:"path"`
	PrettyPath string `json:"pretty_path"`
	NodeIndex  int    `json:"node_index"`
	Url        string `json:"url"`
}
