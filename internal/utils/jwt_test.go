package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateJWT(t *testing.T) {
	t.Parallel()
	var (
		role   = "employee"
		jwtGen = JWTGen{}
	)

	t.Run("no JWT_SECRET", func(t *testing.T) {
		t.Parallel()

		_, err := jwtGen.GenerateJWT(role)

		require.Error(t, err)
	})
}

func Test_GenerateJWTWithSetEnv(t *testing.T) {
	var (
		jwtGen = JWTGen{}
	)
	t.Run("success", func(t *testing.T) {
		role := "employee"
		t.Setenv("JWT_SECRET", "secret")

		_, err := jwtGen.GenerateJWT(role)

		require.NoError(t, err)
	})
}
