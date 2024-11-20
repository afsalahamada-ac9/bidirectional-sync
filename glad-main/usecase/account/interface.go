/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package account

import (
	"sudhagar/glad/entity"
)

// Reader interface
type Reader interface {
	Get(username string) (*entity.Account, error)
	List(tenantID entity.ID) ([]*entity.Account, error)
	GetCount(tenantId entity.ID) (int, error)
}

// Writer interface
type Writer interface {
	Create(e *entity.Account) error
	Update(e *entity.Account) error
	Delete(username string) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	CreateAccount(
		tenantID entity.ID,
		extID string,
		username string,
		first_name string,
		last_name string,
		phone string,
		email string,
		at entity.AccountType) error
	GetAccount(username string) (*entity.Account, error)
	ListAccounts(tenantID entity.ID) ([]*entity.Account, error)
	UpdateAccount(e *entity.Account) error
	DeleteAccount(username string) error
	GetCount(tenantId entity.ID) int
}
