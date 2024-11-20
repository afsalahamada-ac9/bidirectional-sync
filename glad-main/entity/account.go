/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package entity

import (
	"time"
)

// Account type
type AccountType int

const (
	AccountTeacher AccountType = iota
	AccountAssistantTeacher
	AccountOrganizer
	AccountMember
	AccountUser
	// Add new types here
)

// Account data
type Account struct {
	ID ID
	// TenantID ID
	ExtID string

	Username  string
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Type      AccountType

	// meta data
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAccount create a new account
func NewAccount(tenantID ID,
	extID string,
	username string,
	first_name string,
	last_name string,
	phone string,
	email string,
	at AccountType,
) (*Account, error) {
	t := &Account{
		ID:        NewID(),
		ExtID:     extID,
		Username:  username,
		FirstName: first_name,
		LastName:  last_name,
		Phone:     phone,
		Email:     email,
		Type:      at,
		CreatedAt: time.Now(),
	}
	err := t.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return t, nil
}

// Validate validate account
func (t *Account) Validate() error {
	if t.Username == "" {
		return ErrInvalidEntity
	}

	return nil
}
