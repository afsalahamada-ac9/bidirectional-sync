/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package entity

import (
	"time"
)

// Center mode
type CenterMode string

const (
	CenterInPerson CenterMode = "in-person"
	CenterOnline   CenterMode = "online"
	// Add new types here
)

// Center Location
type CenterLocation struct {
	Street1 string
	Street2 string
	City    string
	State   string
	Zip     string
	Country string
}

// Center GeoLocation
type CenterGeoLocation struct {
	Lat  float64
	Long float64
}

// Center data
type Center struct {
	ID       ID
	TenantID ID
	ExtID    string

	Name        string
	Location    CenterLocation
	GeoLocation CenterGeoLocation

	Capacity int32
	Mode     CenterMode

	WebPage          string
	IsNationalCenter bool

	// meta data
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCenterLocation creates a new center location
func NewCenterLocation(street1 string,
	street2 string,
	city string,
	state string,
	zip string,
	country string) (*CenterLocation, error) {

	l := &CenterLocation{
		Street1: street1,
		Street2: street2,
		City:    city,
		State:   state,
		Zip:     zip,
		Country: country,
	}
	err := l.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return l, nil
}

// Validate validates center location
func (l *CenterLocation) Validate() error {
	if l.Street1 == "" || l.City == "" || l.State == "" || l.Zip == "" || l.Country == "" {
		return ErrInvalidEntity
	}
	return nil
}

// NewCenterGeoLocation creates a new center geo location
func NewCenterGeoLocation(lat float64, long float64) (*CenterGeoLocation, error) {
	g := &CenterGeoLocation{
		Lat:  lat,
		Long: long,
	}
	err := g.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return g, nil
}

// Validate validates center geo location
func (g *CenterGeoLocation) Validate() error {
	if g.Lat == 0 || g.Long == 0 {
		return ErrInvalidEntity
	}
	return nil
}

// NewCenter create a new center
func NewCenter(tenantID ID,
	extID string,
	name string,
	location CenterLocation,
	geoLocation CenterGeoLocation,
	capacity int32,
	mode CenterMode,
	webPage string,
	isNationalCenter bool) (*Center, error) {
	c := &Center{
		ID:               NewID(),
		TenantID:         tenantID,
		ExtID:            extID,
		Name:             name,
		Location:         location,
		GeoLocation:      geoLocation,
		Capacity:         capacity,
		Mode:             mode,
		WebPage:          webPage,
		IsNationalCenter: isNationalCenter,
		CreatedAt:        time.Now(),
	}
	err := c.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return c, nil
}

// Validate validate center
func (c *Center) Validate() error {
	if c.ExtID == "" || c.Name == "" {
		return ErrInvalidEntity
	}
	return nil
}
