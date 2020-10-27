package models

import (
	"time"
)

type Check struct {
	ID          int64
	LastUpdated time.Time
	Cloud       string
}

type Service struct {
	Check_id int64
	Name     string
}

type Cloud struct {
	Name              string
	LastUpdated       time.Time
	UnhealthyServices []string
}

type StatusResponse struct {
	Clouds []Cloud
}
