package entity

import "time"

type Tenant struct {
	Id        int
	Name      string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
