package coinstats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBTCDominance(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.BTCDominance(ctx, ChartType1y)
	require.Nil(t, err, err)
	assert.True(t, len(res) > 0)
}

func TestFearAndGreed(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.FearAndGreed(ctx)
	require.Nil(t, err, err)
	assert.True(t, res.Now.Value > 0)
}

func TestFearAndGreedChart(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.FearAndGreedChart(ctx)
	require.Nil(t, err, err)
	assert.True(t, len(res.Data) > 0)
	assert.True(t, res.Data[0].Value > 0)
}

func TestRainbowChart(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.RainbowChart(ctx, "bitcoin")
	require.Nil(t, err, err)
	assert.True(t, len(res) > 0)
	assert.False(t, res[0].Time.IsZero())
}
