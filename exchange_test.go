package coinstats

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExchanges(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.Exchanges(ctx)
	require.Nil(t, err, err)
	assert.Equal(t, res[0].ConnectionID, "gateionative")
	assert.Len(t, res, 47)
}

func TestExchangeBalance(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	_, err := c.ExchangeBalance(ctx, "rULArJ12y+jOQkMtBLekEHV0x+I5pLRhinMGCPzLkRw=", "key", "secret",
		map[string]string{"field": "field", "more": "more"},
	)
	require.NotNil(t, err, err)
	require.True(t, errors.Is(err, ErrBadRequest))
	// !?!?! так АПИ отдает.
	require.True(t, strings.Contains(err.Error(), "connectionFields should not be empty"))
}

func TestExchangeStatus(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	_, err := c.ExchangeStatus(ctx, "unknown")
	require.NotNil(t, err, err)
	// так АПИ отдает. По логике должно быть ErrNotFound. См. тест ниже
	require.True(t, errors.Is(err, ErrBadRequest))
	require.True(t, strings.Contains(err.Error(), "Portfolio not found"))
}

func TestExchangeTransactions(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := ExchangeTransactionFilter{
		Limit: 2,
	}
	_, err := c.ExchangeTransactions(ctx, "unknown", filter)
	require.NotNil(t, err, err)
	require.True(t, errors.Is(err, ErrNotFound))
	require.True(t, strings.Contains(err.Error(), "Portfolio not found"))
}

func TestExchangeChart(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	_, err := c.ExchangeChart(ctx, "unknown", ChartType1w)
	require.NotNil(t, err, err)
	require.True(t, errors.Is(err, ErrNotFound))
	require.True(t, strings.Contains(err.Error(), "Portfolio not found"))
}

func TestExchangeSync(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	_, err := c.ExchangeSync(ctx, "unknown")
	require.NotNil(t, err, err)
	// так АПИ отдает. По логике должно быть ErrNotFound. См. тест выше
	require.True(t, errors.Is(err, ErrBadRequest))
	require.True(t, strings.Contains(err.Error(), "Portfolio not found"))
}
