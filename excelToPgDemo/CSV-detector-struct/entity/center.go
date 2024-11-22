package entity

import "time"

type JSONB string

type Center struct {
	Id               int
	ExtId            string
	CenterName       string
	Location         JSONB // JSONB
	GeoLocation      JSONB // JSONB
	Capacity         int   // INTEGER
	Mode             string
	IsNationalCenter bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
