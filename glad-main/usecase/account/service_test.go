/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package account

import (
	"testing"
	"time"

	"sudhagar/glad/entity"

	"github.com/stretchr/testify/assert"
)

const (
	accountUsernameAlice  string = "12345550001"
	accountUsername2Alice string = "12345550002"

	accountIDAlice  entity.ID = 13790492210917010000
	accountID2Alice entity.ID = 13790492210917010002
	tenantAlice     entity.ID = 13790492210917015554
	aliceExtID                = "001aliceExtID"
	alice2ExtID               = "002aliceExtID"

	// todo: add multi-tenant support
	// tenantBob entity.ID = 13790492210917015555
)

func newFixtureAccount() *entity.Account {
	return &entity.Account{
		ID: accountIDAlice,
		// TenantID:  tenantAlice,
		ExtID:     aliceExtID,
		Username:  accountUsernameAlice,
		FirstName: "Alice",
		LastName:  "Wonderland",
		Phone:     "1235556789",
		Email:     "alice@wonderland.ai",
		Type:      entity.AccountTeacher,
		CreatedAt: time.Now(),
	}
}

func Test_Create(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	account := newFixtureAccount()
	err := m.CreateAccount(tenantAlice,
		account.ExtID,
		account.Username,
		account.FirstName,
		account.LastName,
		account.Phone,
		account.Email,
		account.Type,
	)
	assert.Nil(t, err)
	assert.False(t, account.CreatedAt.IsZero())
}

func Test_SearchAndFind(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	account1 := newFixtureAccount()
	account2 := newFixtureAccount()
	account2.ID = accountID2Alice
	account2.Username = accountUsername2Alice
	account2.ExtID = alice2ExtID

	_ = m.CreateAccount(tenantAlice,
		account1.ExtID,
		account1.Username,
		account1.FirstName,
		account1.LastName,
		account1.Phone,
		account1.Email,
		account1.Type,
	)
	_ = m.CreateAccount(tenantAlice,
		account2.ExtID,
		account2.Username,
		account2.FirstName,
		account2.LastName,
		account2.Phone,
		account2.Email,
		account2.Type,
	)

	t.Run("list all", func(t *testing.T) {
		all, err := m.ListAccounts(tenantAlice)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(all))
	})

	t.Run("get", func(t *testing.T) {
		saved, err := m.GetAccount(account1.Username)
		assert.Nil(t, err)
		assert.Equal(t, account1.ExtID, saved.ExtID)
		assert.Equal(t, account1.Type, saved.Type)
		assert.Equal(t, account1.Username, saved.Username)
	})
}

// It's unlikely that the update will be called in this entity model.
// Perhaps a human readable name can be given for customer to reference.
func Test_Update(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	account := newFixtureAccount()
	err := m.CreateAccount(tenantAlice,
		account.ExtID,
		account.Username,
		account.FirstName,
		account.LastName,
		account.Phone,
		account.Email,
		account.Type,
	)
	assert.Nil(t, err)

	saved, _ := m.GetAccount(account.Username)
	saved.Username = "starred"
	assert.Nil(t, m.UpdateAccount(saved))

	_, err = m.GetAccount(account.Username)
	assert.Equal(t, entity.ErrNotFound, err)

	updated, err := m.GetAccount(saved.Username)
	assert.Nil(t, err)
	assert.Equal(t, "starred", updated.Username)
}

func TestDelete(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)

	account1 := newFixtureAccount()

	account2 := newFixtureAccount()
	account2.ID = accountID2Alice
	account2.Username = accountUsername2Alice
	account2.ExtID = alice2ExtID
	_ = m.CreateAccount(tenantAlice,
		account2.ExtID,
		account2.Username,
		account2.FirstName,
		account2.LastName,
		account2.Phone,
		account2.Email,
		account2.Type,
	)

	err := m.DeleteAccount(account1.Username)
	assert.Equal(t, entity.ErrNotFound, err)

	err = m.DeleteAccount(account2.Username)
	assert.Nil(t, err)

	_, err = m.GetAccount(account2.Username)
	assert.Equal(t, entity.ErrNotFound, err)
}
