package coinstats

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	c := New("")
	_, err := c.CreditUsage(ctx)
	require.NotNil(t, err)
	var e *Error
	require.True(t, errors.As(err, &e))
	require.Equal(t, e.StatusCode, 401)
}
