/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package account

import (
	"time"

	"sudhagar/glad/entity"
)

// Service account usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// CreateAccount create a account
func (s *Service) CreateAccount(
	tenantID entity.ID,
	extID string,
	username string,
	first_name string,
	last_name string,
	phone string,
	email string,
	at entity.AccountType,
) error {
	account, err := entity.NewAccount(tenantID,
		extID,
		username,
		first_name,
		last_name,
		phone,
		email,
		at)
	if err != nil {
		return err
	}
	return s.repo.Create(account)
}

// GetAccount get a account
func (s *Service) GetAccount(username string) (*entity.Account, error) {
	account, err := s.repo.Get(username)
	if account == nil {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

// ListAccounts list account
func (s *Service) ListAccounts(tenantID entity.ID) ([]*entity.Account, error) {
	accounts, err := s.repo.List(tenantID)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, entity.ErrNotFound
	}
	return accounts, nil
}

// UpdateAccount Update a account
func (s *Service) UpdateAccount(t *entity.Account) error {
	err := t.Validate()
	if err != nil {
		return err
	}
	t.UpdatedAt = time.Now()
	return s.repo.Update(t)
}

// DeleteAccount Delete a account
func (s *Service) DeleteAccount(username string) error {
	account, err := s.GetAccount(username)
	if account == nil {
		return entity.ErrNotFound
	}
	if err != nil {
		return err
	}

	return s.repo.Delete(username)
}

// GetCount gets total account count
func (s *Service) GetCount(tenantID entity.ID) int {
	count, err := s.repo.GetCount(tenantID)
	if err != nil {
		return 0
	}

	return count
}
