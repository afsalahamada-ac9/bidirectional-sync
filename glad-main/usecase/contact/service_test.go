/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package contact

import (
	"testing"

	"sudhagar/glad/entity"

	"github.com/stretchr/testify/assert"
)

const (
	contactHandleAlice string = "14565550001"
	contactHandleBob   string = "14565550002"

	contactIDAlice entity.ID = 13790492210917010000
	contactIDBob   entity.ID = 13790492210917010002

	nameAlice = "Alice Wonderland"
	nameBob   = "Bob Builder"

	tenantAlpha entity.ID = 13790492210917015554

	accountAlphaPrimary entity.ID = 13790492210917011000

	// todo: add multi-account support
	// accountAlphaSecondary entity.ID = 13790492210917011001

	// todo: add multi-tenant support
	// tenantBeta entity.ID = 13790492210917015555
)

func newFixtureContactAlice() *entity.Contact {
	return &entity.Contact{
		ID:        contactIDAlice,
		TenantID:  tenantAlpha,
		AccountID: accountAlphaPrimary,
		Handle:    contactHandleAlice,
		Name:      nameAlice,
		IsStale:   false,
	}
}

func newFixtureContactBob() *entity.Contact {
	return &entity.Contact{
		ID:        contactIDBob,
		TenantID:  tenantAlpha,
		AccountID: accountAlphaPrimary,
		Handle:    contactHandleBob,
		Name:      nameBob,
		IsStale:   false,
	}
}

func Test_CreateOrUpdate(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	contact := newFixtureContactAlice()

	t.Run("Create", func(t *testing.T) {
		err := m.CreateOrUpdate(contact.TenantID, contact.AccountID, contact.Handle, contact.Name)
		assert.Nil(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		err := m.CreateOrUpdate(contact.TenantID, contact.AccountID, contact.Handle, "Alice Wunderland")
		assert.Nil(t, err)
	})
}

func Test_StaleByAccount(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	contactA := newFixtureContactAlice()
	err := m.CreateOrUpdate(contactA.TenantID, contactA.AccountID, contactA.Handle, contactA.Name)
	assert.Nil(t, err)

	contactB := newFixtureContactBob()
	err = m.CreateOrUpdate(contactB.TenantID, contactB.AccountID, contactB.Handle, contactB.Name)
	assert.Nil(t, err)

	t.Run("Set stale", func(t *testing.T) {
		err = m.SetStaleByAccount(contactA.TenantID, contactA.AccountID)
		assert.Nil(t, err)

		contacts, err := m.GetMulti(contactA.TenantID, 0, 10)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(contacts))
		assert.Equal(t, true, contacts[0].IsStale)
		assert.Equal(t, true, contacts[1].IsStale)
	})

	t.Run("Reset stale", func(t *testing.T) {
		err = m.ResetStaleByAccount(contactA.TenantID, contactA.AccountID)
		assert.Nil(t, err)

		contacts, err := m.GetMulti(contactA.TenantID, 0, 10)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(contacts))
		assert.Equal(t, false, contacts[0].IsStale)
		assert.Equal(t, false, contacts[1].IsStale)
	})

	t.Run("Delete stale", func(t *testing.T) {
		err = m.DeleteStaleByAccount(contactA.TenantID, contactA.AccountID)
		assert.Nil(t, err)

		contacts, err := m.GetMulti(contactA.TenantID, 0, 10)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(contacts))

		err = m.SetStaleByAccount(contactA.TenantID, contactA.AccountID)
		assert.Nil(t, err)
		err = m.DeleteStaleByAccount(contactA.TenantID, contactA.AccountID)
		assert.Nil(t, err)

		contacts, err = m.GetMulti(contactA.TenantID, 0, 10)
		assert.Equal(t, entity.ErrNotFound, err)
		assert.Nil(t, contacts)
	})
}

func TestDelete(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	contactA := newFixtureContactAlice()
	err := m.CreateOrUpdate(contactA.TenantID, contactA.AccountID, contactA.Handle, contactA.Name)
	assert.Nil(t, err)

	contactB := newFixtureContactBob()
	err = m.CreateOrUpdate(contactB.TenantID, contactB.AccountID, contactB.Handle, contactB.Name)
	assert.Nil(t, err)

	contacts, err := m.GetMulti(contactA.TenantID, 0, 10)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(contacts))

	err = m.DeleteByAccount(contactA.TenantID, contactA.AccountID)
	assert.Nil(t, err)
}

func TestGetCount(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	contactA := newFixtureContactAlice()
	err := m.CreateOrUpdate(contactA.TenantID, contactA.AccountID, contactA.Handle, contactA.Name)
	assert.Nil(t, err)

	contactB := newFixtureContactBob()
	err = m.CreateOrUpdate(contactB.TenantID, contactB.AccountID, contactB.Handle, contactB.Name)
	assert.Nil(t, err)

	count, err := m.GetCount(contactA.TenantID)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}
