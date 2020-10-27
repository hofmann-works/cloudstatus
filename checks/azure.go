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

	/*if azresponse.Status.Health == "healthy" {
		check := &models.Check{Cloud: "Azure", LastUpdated: azresponse.LastUpdated}
		database.AddCheck(check)
		fmt.Println("ID:", check.ID)
	}*/
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
