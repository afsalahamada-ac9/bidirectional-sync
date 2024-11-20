/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package label

import (
	"testing"

	"sudhagar/glad/entity"

	"github.com/stretchr/testify/assert"
)

const (
	labelIDAlice entity.ID = 13790492210917020000
	labelIDBob   entity.ID = 13790492210917020002

	nameAlice string = "Alice"
	nameBob   string = "Bob"
	nameNew   string = "New"

	tenantAlpha entity.ID = 13790492210917015554

	colorAlice uint32 = 111
	colorBob   uint32 = 222

	// todo: add multi-tenant support
	// tenantBeta entity.ID = 13790492210917015555
)

func newFixtureLabelAlice() *entity.Label {
	return &entity.Label{
		ID:       labelIDAlice,
		TenantID: tenantAlpha,
		Name:     nameAlice,
		Color:    colorAlice,
	}
}

func newFixtureLabelBob() *entity.Label {
	return &entity.Label{
		ID:       labelIDBob,
		TenantID: tenantAlpha,
		Name:     nameBob,
		Color:    colorBob,
	}
}

func Test_Create(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	label := newFixtureLabelAlice()

	_, err := m.Create(label.TenantID, label.Name, label.Color)
	assert.Nil(t, err)
}

func Test_Update(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	label := newFixtureLabelAlice()

	err := m.Update(label)
	assert.Equal(t, entity.ErrNotFound, err)

	labelID, err := m.Create(label.TenantID, label.Name, label.Color)
	assert.Nil(t, err)

	label.ID = labelID
	label.Name = nameNew
	err = m.Update(label)
	assert.Nil(t, err)

	got, err := m.Get(label.TenantID, labelID)
	assert.Nil(t, err)
	assert.Equal(t, label.ID, got.ID)
	assert.Equal(t, nameNew, got.Name)
	assert.Equal(t, label.Color, got.Color)
}

func TestDelete(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	labelA := newFixtureLabelAlice()

	labelAID, err := m.Create(labelA.TenantID, labelA.Name, labelA.Color)
	assert.Nil(t, err)

	labelB := newFixtureLabelBob()
	_, err = m.Create(labelB.TenantID, labelB.Name, labelB.Color)
	assert.Nil(t, err)

	labels, err := m.GetMulti(labelA.TenantID, 0, 10)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(labels))

	err = m.Delete(labelA.TenantID, labelAID)
	assert.Nil(t, err)

	_, err = m.Get(labelA.TenantID, labelAID)
	assert.Equal(t, entity.ErrNotFound, err)
}

func TestGetCount(t *testing.T) {
	repo := newInmem()
	m := NewService(repo)
	labelA := newFixtureLabelAlice()
	_, err := m.Create(labelA.TenantID, labelA.Name, labelA.Color)
	assert.Nil(t, err)

	labelB := newFixtureLabelBob()
	_, err = m.Create(labelB.TenantID, labelB.Name, labelB.Color)
	assert.Nil(t, err)

	count, err := m.GetCount(labelA.TenantID)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}
