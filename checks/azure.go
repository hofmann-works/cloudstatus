package checks

import (
	"encoding/json"
	"fmt"
	"github.com/hofmann-works/cloudstatus/config"
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

func AzureStatus() {
	azureStatusURL := config.New().AzureStatusURL

	response, err := http.Get(azureStatusURL)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var azresponse AzResponse
		err := json.Unmarshal([]byte(data), &azresponse)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(azresponse.LastUpdated, azresponse.Services[0].Geographies[1].Name)
	}
}
