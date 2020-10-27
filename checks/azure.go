package checks

import (
	"encoding/json"
	"github.com/hofmann-works/cloudstatus/config"
	"github.com/hofmann-works/cloudstatus/db"
	"github.com/hofmann-works/cloudstatus/models"
	"io/ioutil"
	"log"
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
		log.Printf("Azure HTTP request failed with error %s\n", err)
		return
	}

	data, _ := ioutil.ReadAll(response.Body)

	//fakedata
	//data = []byte(`{"lastUpdated":"2020-10-21T07:15:04.703Z","status":{"health":"healthy","message":"Everything is looking good"},"services":[{"id":"Core services","geographies":[{"id":"US","name":"United States","health":"healthy"},{"id":"EU","name":"Europe","health":"healthy"}]},{"id":"Boards","geographies":[{"id":"US","name":"United States","health":"healthy"},{"id":"EU","name":"Europe","health":"healthy"}]},{"id":"Repos","geographies":[{"id":"US","name":"United States","health":"healthy"},{"id":"EU","name":"Europe","health":"unhealthy"}]},{"id":"Pipelines","geographies":[{"id":"US","name":"United States","health":"unhealthy"},{"id":"EU","name":"Europe","health":"healthy"}]},{"id":"Test Plans","geographies":[{"id":"US","name":"United States","health":"healthy"},{"id":"EU","name":"Europe","health":"healthy"}]},{"id":"Artifacts","geographies":[{"id":"US","name":"United States","health":"healthy"},{"id":"EU","name":"Europe","health":"healthy"}]},{"id":"Other services","geographies":[{"id":"US","name":"United States","health":"healthy"},{"id":"EU","name":"Europe","health":"healthy"}]}]}`)

	err = json.Unmarshal([]byte(data), &azresponse)
	if err != nil {
		log.Printf("Unmarshal Azure HTTP response failed with error %s\n", err)
	}

	serviceIsUnhealthy := false
	unhelathyServices := []string{}

	for _, service := range azresponse.Services {

		for _, geopgraphy := range service.Geographies {
			if geopgraphy.Health != "healthy" {
				serviceIsUnhealthy = true
			}
		}

		if serviceIsUnhealthy {
			unhelathyServices = append(unhelathyServices, service.Id)
			serviceIsUnhealthy = false
		}
	}

	check := &models.Check{Cloud: "Azure", LastUpdated: azresponse.LastUpdated}
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
