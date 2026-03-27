package coinstats

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNftTrending(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := NftTrendingFilter{Limit: 5}
	res, err := c.NftTrending(ctx, filter)
	require.Nil(t, err, err)
	assert.Equal(t, res.Meta.Limit, 5)
	assert.Equal(t, res.Meta.ItemCount, 223016)
	assert.Len(t, res.Data, 5)
}

func TestNftsByWallet(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := NftsByWalletFilter{Limit: 5}
	res, err := c.NftsByWallet(ctx, "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", filter)
	require.Nil(t, err, err)
	assert.Equal(t, res.Meta.Limit, 5)
	assert.Equal(t, res.Meta.ItemCount, 11)
	assert.Len(t, res.Data, 5)
}

func TestNftCollectionByAddress(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.NftCollectionByAddress(ctx, "0x139fcbe60644dfc74f80f7c392beb6d1783f53bd")
	require.Nil(t, err, err)
	assert.Equal(t, res.Address, "0x139fcbe60644dfc74f80f7c392beb6d1783f53bd")
}

func TestNftCollectionAssetsByAddress(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := NftCollectionAssetFilter{Limit: 5}
	res, err := c.NftCollectionAssetsByAddress(ctx, "0x139fcbe60644dfc74f80f7c392beb6d1783f53bd", filter)
	require.Nil(t, err, err)
	assert.Equal(t, res.Meta.Limit, 5)
	assert.Equal(t, res.Meta.ItemCount, 1)
	assert.Len(t, res.Data, 1)
	fmt.Printf("%+v\n", res.Data)
}

func TestNftCollectionAssetByToken(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.NftCollectionAssetByToken(ctx, "0x139fcbe60644dfc74f80f7c392beb6d1783f53bd", "0")
	require.Nil(t, err, err)
	assert.Equal(t, res.Address, "0x139fcbe60644dfc74f80f7c392beb6d1783f53bd")
	assert.Equal(t, res.TokenID, "0")
}
