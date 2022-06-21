package apperrors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thedhejavu/errauditor/examples/project/pkg/apperrors"
)

func TestComparingErrors(t *testing.T) {
	t.Parallel()

	t.Run("Dynamic domain errors are comparable if they have the same code", func(t *testing.T) {
		t.Parallel()

		err1 := apperrors.NewDomainError("code", "message")
		err2 := apperrors.NewDomainError("code", "message2")

		require.True(t, errors.Is(err1, err2))

		err3 := apperrors.NewDomainError("code2", "message")
		require.False(t, errors.Is(err1, err3))
	})

	t.Run("Static app errors are comparable", func(t *testing.T) {
		t.Parallel()

		err := apperrors.ErrAddressNotFound
		require.True(t, errors.Is(err, apperrors.ErrAddressNotFound))
	})
}
