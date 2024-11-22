package entity

import "time"

type Account struct {
	Id        int
	ExtId     string
	Username  string
	FirstName string // VARCHAR(40)
	LastName  string // VARCHAR(80)
	Phone     string // VARCHAR(32)
	Email     string // VARCHAR(80)
	Type      string // ENUM
	CreatedAt time.Time
	UpdatedAt time.Time
}
