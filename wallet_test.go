package coinstats

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockchains(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.Blockchains(ctx)
	require.Nil(t, err, err)
	assert.True(t, len(res) > 0)
	assert.Equal(t, res[0].Chain, "binance_smart")
}

func TestWalletBalance(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	wallet := Wallet{
		Address: "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
	}
	_, err := c.WalletBalance(ctx, wallet)
	require.NotNil(t, err)
	wallet.Blockchain = "base"
	res, err := c.WalletBalance(ctx, wallet)
	require.NotNil(t, err)
	assert.Len(t, res, 0)
}

func TestWalletBalances(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := WalletBalancesFilter{}
	_, err := c.WalletBalances(ctx, filter)
	require.NotNil(t, err, err)

	filter.Address = "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	filter.Blockchains = []string{"ethereum"}
	res, err := c.WalletBalances(ctx, filter)
	require.Nil(t, err, err)
	assert.True(t, len(res) > 0)
	assert.Equal(t, res[0].Blockchain, "ethereum")
}

func TestWalletSyncStatus(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	wallet := Wallet{
		Address:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		Blockchain: "base",
	}
	res, err := c.WalletSyncStatus(ctx, wallet)
	require.Nil(t, err, err)
	assert.Equal(t, res, SyncedStatus)
}

func TestWalletTransactions(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := WalletTransactionFilter{
		Wallet: Wallet{
			Address:      "0x1234567890abcdef1234567890abcdef12345678",
			ConnectionID: "ethereum",
		},
	}
	res, err := c.WalletTransactions(ctx, filter)
	require.Nil(t, err, err)
	assert.Len(t, res.Result, 20)

	filter.DateFrom = time.Date(2025, 4, 6, 20, 45, 0, 0, time.UTC)
	res, err = c.WalletTransactions(ctx, filter)
	require.Nil(t, err, err)
	assert.Len(t, res.Result, 9)
}

func TestWalletTransactionsSync(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := WalletTransactionsSyncFilter{
		Wallet: Wallet{
			Address:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			Blockchain: "ethereum",
		},
	}
	res, err := c.WalletTransactionsSync(ctx, filter)
	require.Nil(t, err)
	require.Equal(t, res, SyncingStatus)

	filter.Wallets = []WalletShort{
		{
			Address:      "0x1234567890abcdef1234567890abcdef12345678",
			ConnectionID: "ethereum",
		},
		{
			Address:      "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			ConnectionID: "ethereum",
		},
	}
	res, err = c.WalletTransactionsSync(ctx, filter)
	require.Nil(t, err, err)
	require.Equal(t, res, SyncingStatus)
}

func TestWalletChart(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.WalletChart(ctx, ChartType1y, Wallet{
		Address:      "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		ConnectionID: "ethereum",
	})
	require.Nil(t, err, err)
	require.True(t, len(res) > 0)
}

func TestWalletCharts(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.WalletCharts(ctx, ChartType1y, []WalletShort{
		{
			Address:      "0x1234567890abcdef1234567890abcdef12345678",
			ConnectionID: "ethereum",
		},
		{
			Address:      "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			ConnectionID: "ethereum",
		},
	}, false)
	require.Nil(t, err, err)
	require.Len(t, res, 2)
	require.True(t, len(res[0].Data) > 0)
	require.Equal(t, res[0].ConnectionID, "ethereum")
}

func TestWalletDefi(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	w := Wallet{
		Address:      "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		ConnectionID: "ethereum",
	}
	res, err := c.WalletDefi(ctx, w)
	require.Nil(t, err, err)
	require.Greater(t, res.TotalAssets.BTC, 0.0)
	require.Greater(t, res.TotalAssets.ETH, 0.0)
	require.Greater(t, res.TotalAssets.USD, 0.0)
}
