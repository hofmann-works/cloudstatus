package checks

import (
	"encoding/json"
	"fmt"
	"github.com/hofmann-works/cloudstatus/config"
	"io/ioutil"
	"net/http"
	"time"
)

type GitHubResponse struct {
	Page struct {
		Id         string
		Name       string
		Url        string
		Time_zone  string
		Updated_at time.Time
	}
	Components []struct {
		Id     string
		Name   string
		Status string
	}
}

func GitHubStatus() {
	GitHubStatusURL := config.New().GitHubStatusURL
	response, err := http.Get(GitHubStatusURL)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var githubresponse GitHubResponse
		err := json.Unmarshal([]byte(data), &githubresponse)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(githubresponse.Page.Updated_at, githubresponse.Components[1].Name)
	}
}
