package coinstats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreditUsage(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.CreditUsage(ctx)
	require.Nil(t, err, err)
	assert.Greater(t, res.Total, 0)
	assert.Greater(t, res.Used, 0)
	assert.Greater(t, res.Remaining, 0)
}
