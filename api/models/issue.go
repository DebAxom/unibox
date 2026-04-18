package models

import "time"

type Issue struct {
	Id         string
	Issuer     string
	Title      string
	Desc       string
	Img        string // image url only
	Status     string
	Dept       string
	Updated_at time.Time
}
