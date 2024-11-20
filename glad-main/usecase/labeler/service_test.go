/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package labeler

import (
	"testing"

	"sudhagar/glad/entity"

	cmock "sudhagar/glad/usecase/contact/mock"
	lmock "sudhagar/glad/usecase/label/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	tenantAlpha entity.ID = 13790492210917015554

	contactIDAlice     entity.ID = 13790492210917010000
	contactHandleAlice string    = "14565550001"
	contactNameAlice             = "Alice Wonderland"

	labelIDAlice entity.ID = 13790492210917020000

	// todo: add multi-tenant support
	// tenantBeta entity.ID = 13790492210917015555
)

func newContactAlice() *entity.Contact {
	return &entity.Contact{
		ID:       contactIDAlice,
		TenantID: tenantAlpha,
		Handle:   contactHandleAlice,
		Name:     contactNameAlice,
		// AccountID, Handle, Name, IsStale are not required
	}
}

func newContactAliceWithLabel(labelID entity.ID) *entity.Contact {
	c := newContactAlice()
	c.Labels = append(c.Labels, labelID)
	return c
}

func newLabelAlice() *entity.Label {
	return &entity.Label{
		ID:       labelIDAlice,
		TenantID: tenantAlpha,
		// Name & Color are not required
	}
}

func Test_SetLabel(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	cMock := cmock.NewMockUseCase(controller)
	lMock := lmock.NewMockUseCase(controller)
	uc := NewService(cMock, lMock)

	cMock.EXPECT().Update(gomock.Any()).Return(nil)
	lMock.EXPECT().Update(gomock.Any()).Return(nil)

	t.Run("Happy Path", func(t *testing.T) {
		cMock.EXPECT().Get(tenantAlpha, contactIDAlice).Return(newContactAlice(), nil)
		lMock.EXPECT().Get(tenantAlpha, labelIDAlice).Return(newLabelAlice(), nil)

		err := uc.SetLabel(tenantAlpha, contactIDAlice, labelIDAlice)
		assert.Nil(t, err)
	})

	t.Run("Contact not found", func(t *testing.T) {
		cMock.EXPECT().Get(tenantAlpha, contactIDAlice).Return(nil, entity.ErrNotFound)
		err := uc.SetLabel(tenantAlpha, contactIDAlice, labelIDAlice)
		assert.Equal(t, entity.ErrNotFound, err)
	})

	t.Run("Label not found", func(t *testing.T) {
		cMock.EXPECT().Get(tenantAlpha, contactIDAlice).Return(newContactAlice(), nil)
		lMock.EXPECT().Get(tenantAlpha, labelIDAlice).Return(nil, entity.ErrNotFound)
		err := uc.SetLabel(tenantAlpha, contactIDAlice, labelIDAlice)
		assert.Equal(t, entity.ErrNotFound, err)
	})

	t.Run("Label already set", func(t *testing.T) {
		cMock.EXPECT().Get(tenantAlpha, contactIDAlice).Return(newContactAliceWithLabel(labelIDAlice), nil)
		lMock.EXPECT().Get(tenantAlpha, labelIDAlice).Return(newLabelAlice(), nil)
		err := uc.SetLabel(tenantAlpha, contactIDAlice, labelIDAlice)
		assert.Equal(t, entity.ErrLabelAlreadySet, err)
	})

	// Too many labels case must be tested in the entity -> contact_test.go
}

func Test_RemoveLabel(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	cMock := cmock.NewMockUseCase(controller)
	lMock := lmock.NewMockUseCase(controller)
	uc := NewService(cMock, lMock)

	cMock.EXPECT().Update(gomock.Any()).Return(nil)
	lMock.EXPECT().Update(gomock.Any()).Return(nil)

	t.Run("Happy Path", func(t *testing.T) {
		cMock.EXPECT().Get(tenantAlpha, contactIDAlice).Return(newContactAliceWithLabel(labelIDAlice), nil)
		lMock.EXPECT().Get(tenantAlpha, labelIDAlice).Return(newLabelAlice(), nil)

		err := uc.RemoveLabel(tenantAlpha, contactIDAlice, labelIDAlice)
		assert.Nil(t, err)
	})

	t.Run("Contact not found", func(t *testing.T) {
		cMock.EXPECT().Get(tenantAlpha, contactIDAlice).Return(nil, entity.ErrNotFound)
		err := uc.RemoveLabel(tenantAlpha, contactIDAlice, labelIDAlice)
		assert.Equal(t, entity.ErrNotFound, err)
	})

	t.Run("Label not found", func(t *testing.T) {
		cMock.EXPECT().Get(tenantAlpha, contactIDAlice).Return(newContactAlice(), nil)
		lMock.EXPECT().Get(tenantAlpha, labelIDAlice).Return(nil, entity.ErrNotFound)

		err := uc.RemoveLabel(tenantAlpha, contactIDAlice, labelIDAlice)
		assert.Equal(t, entity.ErrNotFound, err)
	})

	t.Run("Label not set", func(t *testing.T) {
		cMock.EXPECT().Get(tenantAlpha, contactIDAlice).Return(newContactAlice(), nil)
		lMock.EXPECT().Get(tenantAlpha, labelIDAlice).Return(newLabelAlice(), nil)

		err := uc.RemoveLabel(tenantAlpha, contactIDAlice, labelIDAlice)
		assert.Equal(t, entity.ErrNotFound, err)
	})
}
