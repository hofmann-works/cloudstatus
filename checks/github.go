package checks

import (
	"encoding/json"
	"fmt"
	"github.com/hofmann-works/cloudstatus/config"
	"github.com/hofmann-works/cloudstatus/db"
	"github.com/hofmann-works/cloudstatus/models"
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

func GitHubStatus(database db.Database) {
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

		unhelathyServices := []string{}

		for _, service := range githubresponse.Components {
			if service.Status != "operational" {
				unhelathyServices = append(unhelathyServices, service.Name)
			}
		}

		check := &models.Check{Cloud: "GitHub", LastUpdated: githubresponse.Page.Updated_at}
		service := &models.Service{}
		database.AddCheck(check)
		if check.ID == 0 {
			//check already exist
			return
		} else {
			service.Check_id = check.ID
			for _, unhealthyService := range unhelathyServices {
				service.Name = unhealthyService
				database.AddService(service)
			}
		}
	}
}
