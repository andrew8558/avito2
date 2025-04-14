package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidCity(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		city := CityMoscow

		ok := city.IsValid()

		assert.True(t, ok)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		var city City = "test"

		ok := city.IsValid()

		assert.False(t, ok)
	})
}

func Test_IsValidProductType(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		productType := ProductTypeClothes

		ok := productType.IsValid()

		assert.True(t, ok)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		var productType ProductType = "test"

		ok := productType.IsValid()

		assert.False(t, ok)
	})
}

func Test_IsValidRole(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		role := RoleEmployee

		ok := role.IsValid()

		assert.True(t, ok)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		var role Role = "test"

		ok := role.IsValid()

		assert.False(t, ok)
	})
}
