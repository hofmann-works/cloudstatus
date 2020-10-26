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

type AzResponse struct {
	LastUpdated time.Time
	Status      struct {
		Health  string
		Message string
	}
	Services []struct {
		Id          string
		Geographies []struct {
			Id     string
			Name   string
			Health string
		}
	}
}

func AzureStatus(database db.Database) {
	azureStatusURL := config.New().AzureStatusURL
	var azresponse AzResponse

	response, err := http.Get(azureStatusURL)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	data, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal([]byte(data), &azresponse)
	if err != nil {
		fmt.Println(err)
	}

	if azresponse.Status.Health == "healthy" {
		check := &models.Check{Cloud: "Azure", LastUpdated: azresponse.LastUpdated}
		database.AddCheck(check)
		fmt.Println("ID:", check.ID)
	}

}
