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
