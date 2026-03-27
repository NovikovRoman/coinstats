package coinstats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoins(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := CoinFilter{
		Limit: 50,
	}
	res, err := c.Coins(ctx, filter)
	require.Nil(t, err, err)
	assert.Equal(t, res.Meta.Limit, 50)
	assert.True(t, res.Meta.HasNextPage)
	assert.False(t, res.Meta.HasPreviousPage)
	assert.Len(t, res.Result, 50)
}

func TestCoinCharts(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.CoinCharts(ctx, ChartType1m, []string{"bitcoin", "etherium"})
	require.Nil(t, err, err)
	assert.Equal(t, res[0].CoinID, "bitcoin")
	assert.Len(t, res, 2)
}

func TestCoinByID(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.CoinByID(ctx, "bitcoin", "USD")
	require.Nil(t, err, err)
	assert.Equal(t, res.ID, "bitcoin")
	assert.Equal(t, res.Slug, "bitcoin")
}

func TestCoinChartByID(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.CoinChartByID(ctx, ChartType1m, "bitcoin")
	require.Nil(t, err, err)
	assert.True(t, len(res) > 0)
}

func TestCoinAvgPrice(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.CoinAvgPrice(ctx, "bitcoin", 1739428176)
	require.Nil(t, err, err)
	assert.True(t, res.BTC > 0)
	assert.True(t, res.ETH > 0)
	assert.True(t, res.USD > 0)
}

func TestCoinExchangePrice(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.CoinExchangePrice(ctx, "Binance", "BTC", "ETH", 1739428176)
	require.Nil(t, err, err)
	assert.Equal(t, res, 35.75259206292456)
}

func TestTickerExchanges(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.TickerExchanges(ctx)
	require.Nil(t, err, err)
	assert.Equal(t, res[0].Name, "Poloniex")
	assert.Len(t, res, 270)
}

func TestTickerMarkets(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	filter := TickerMarketsFilter{
		Limit: 10,
		Page:  2,
	}
	res, err := c.TickerMarkets(ctx, filter)
	require.Nil(t, err, err)
	assert.Equal(t, res.Meta.Page, 2)
	assert.Equal(t, res.Meta.Limit, 10)
	assert.Len(t, res.Result, 10)
}

func TestFiats(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.Fiats(ctx)
	require.Nil(t, err, err)
	assert.Equal(t, res[0].Name, "AUD")
	assert.Len(t, res, 63)
}

func TestMarkets(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.Markets(ctx)
	require.Nil(t, err, err)
	assert.True(t, res.Volume > 0)
}

func TestCurrencies(t *testing.T) {
	ctx := context.Background()
	c := New(testApiKey)
	res, err := c.Currencies(ctx)
	require.Nil(t, err, err)
	assert.True(t, res["USD"] > 0)
	assert.True(t, res["CHF"] > 0)
}
