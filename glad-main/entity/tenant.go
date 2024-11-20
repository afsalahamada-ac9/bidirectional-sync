/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package entity

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Tenant data
type Tenant struct {
	ID      ID
	Name    string
	Country string

	AuthToken string

	// meta data
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewTenant create a new tenant
func NewTenant(name, country string) (*Tenant, error) {
	t := &Tenant{
		ID:        NewID(),
		Name:      name,
		Country:   country,
		CreatedAt: time.Now(),
	}
	err := t.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return t, nil
}

// Validate validate tenant
func (t *Tenant) Validate() error {
	if t.Country == "" {
		log.Printf("Invalid country")
		return ErrInvalidEntity
	}

	if t.Name == "" {
		return ErrInvalidEntity
	}
	return nil
}

// UNUSED: ValidatePassword validate tenant password
func (t *Tenant) ValidatePassword(p1, p2 string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p1), []byte(p2))
	if err != nil {
		return err
	}
	return nil
}

// UNUSED
// func genPassword(raw string) string {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
// 	if err != nil {
// 		log.Printf("Unable to generate password, %v", err.Error())
// 		return ""
// 	}

// 	return string(hash)
// }

// ValidateToken validate tenant auth token
func (t *Tenant) ValidateToken(token, password string) error {
	if token != t.genToken(password) {
		return ErrTokenMismatch
	}

	return nil
}

// simple token generator for now
func (t *Tenant) genToken(password string) string {
	token := sha256.New()
	token.Write([]byte(password))
	hash := token.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash)
}

// Exported API to generate and store the token
func (t *Tenant) GenToken(password string) error {
	t.AuthToken = t.genToken(password)
	return nil
}
