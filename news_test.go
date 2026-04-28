package coinstats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewsSources(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.NewsSources(ctx)
	require.Nil(t, err, err)
	assert.Len(t, res, 161)
}

func TestNews(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := NewsFilter{
		Limit: 5,
	}
	res, err := c.News(ctx, filter)
	require.Nil(t, err, err)
	assert.Len(t, res, 5)
}

func TestNewsByType(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := NewsByTypeFilter{
		Limit: 5,
	}
	res, err := c.NewsByType(ctx, NewsTypeBearish, filter)
	require.Nil(t, err, err)
	assert.Len(t, res, 5)
}
func TestNewsByID(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	id := "c5667f721751911cfee05289b408d205225b1928c1aa9f25e4914e4f0e8c1353"
	res, err := c.NewsByID(ctx, id)
	require.Nil(t, err, err)
	assert.Equal(t, res.ID, id)
}
