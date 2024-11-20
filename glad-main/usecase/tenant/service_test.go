/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package tenant

import (
	"testing"
	"time"

	"sudhagar/glad/entity"

	"github.com/stretchr/testify/assert"
)

const (
	tenantDefault entity.ID = 13790492210917015554
)

func newFixtureTenant() *entity.Tenant {
	return &entity.Tenant{
		ID:        tenantDefault,
		Name:      "alice@wonder.land",
		Country:   "testing123",
		AuthToken: "token123",
		CreatedAt: time.Now(),
	}
}

func Test_Create(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	tenant := newFixtureTenant()
	_, err := m.CreateTenant(tenant.Name, tenant.Country)
	assert.Nil(t, err)
	assert.False(t, tenant.CreatedAt.IsZero())
}

func Test_SearchAndFind(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	t1 := newFixtureTenant()
	t2 := newFixtureTenant()
	t2.Name = "bob@wunder.land"

	tID, _ := m.CreateTenant(t1.Name, t1.Country)
	_, _ = m.CreateTenant(t2.Name, t2.Country)

	t.Run("list all", func(t *testing.T) {
		all, err := m.ListTenants()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(all))
	})

	t.Run("get", func(t *testing.T) {
		saved, err := m.GetTenant(tID)
		assert.Nil(t, err)
		// do not compare ID, because it's generated in CreateTenant
		assert.Equal(t, t1.Name, saved.Name)
		// do not compare password, because password is encrypted
		// in CreateTenant
	})
}

func Test_Update(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	tenant := newFixtureTenant()
	id, err := m.CreateTenant(tenant.Name, tenant.Country)
	assert.Nil(t, err)

	saved, _ := m.GetTenant(id)
	saved.Country = "testing456"
	assert.Nil(t, m.UpdateTenant(saved))

	updated, err := m.GetTenant(id)
	assert.Nil(t, err)
	assert.Equal(t, saved.Country, updated.Country)
}

func TestDelete(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)

	t1 := newFixtureTenant()
	t2 := newFixtureTenant()
	t2ID, _ := m.CreateTenant(t2.Name, t2.Country)

	err := m.DeleteTenant(t1.ID)
	assert.Equal(t, entity.ErrNotFound, err)

	err = m.DeleteTenant(t2ID)
	assert.Nil(t, err)
	_, err = m.GetTenant(t2ID)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestLogin(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)

	// pre-requisities
	t1 := newFixtureTenant()
	t2 := newFixtureTenant()
	t2.Name = "bob@wunder.land"

	t1ID, _ := m.CreateTenant(t1.Name, t1.Country)
	_, _ = m.CreateTenant(t2.Name, t2.Country)

	// test
	t.Run("valid credentials", func(t *testing.T) {
		tenant, err := m.Login(t1.Name, t1.Country)
		assert.Nil(t, err)
		assert.Equal(t, t1.Name, tenant.Name)
		assert.Equal(t, t1ID, tenant.ID)
	})

	// Not a valid test case
	// t.Run("invalid credentials", func(t *testing.T) {
	// 	tenant, err := m.Login(t1.Name, "CountryInvalid")
	// 	assert.Equal(t, err, entity.ErrAuthFailure)
	// 	assert.Nil(t, tenant)
	// })
}
