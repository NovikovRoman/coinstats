package coinstats

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CoinFilter struct {
	Page                  int
	Limit                 int
	CoinIDs               []string
	Currency              string
	Name                  string
	Symbol                string
	Blockchains           []string
	IncludeRiskScore      *bool
	Categories            []string
	SortBy                string
	SortDir               SortDir
	MarketCap             Thresholds
	FullyDilutedValuation Thresholds
	Volume                Thresholds
	PriceChange1h         Thresholds
	PriceChange1d         Thresholds
	PriceChange7d         Thresholds
	AvailableSupply       Thresholds
	TotalSupply           Thresholds
	Rank                  Thresholds
	Price                 Thresholds
	RiskScore             Thresholds
}

type CoinResult struct {
	Meta   Meta       `json:"meta"`
	Result []CoinBase `json:"result"`
}

// Coins return comprehensive data about all cryptocurrencies.
func (c *Client) Coins(ctx context.Context, filter CoinFilter) (CoinResult, error) {
	q := url.Values{}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	if len(filter.CoinIDs) > 0 {
		q.Add("coinIds", strings.Join(filter.CoinIDs, ","))
	}
	if filter.Currency != "" {
		q.Add("currency", filter.Currency)
	}
	if filter.Name != "" {
		q.Add("name", filter.Name)
	}
	if filter.Symbol != "" {
		q.Add("symbol", filter.Symbol)
	}
	if len(filter.Blockchains) > 0 {
		q.Add("blockchains", strings.Join(filter.Blockchains, ","))
	}
	if len(filter.Categories) > 0 {
		q.Add("categories", strings.Join(filter.Categories, ","))
	}
	if filter.SortBy != "" {
		q.Add("sortBy", filter.SortBy)
	}
	if filter.SortDir != "" {
		q.Add("sortDir", string(filter.SortDir))
	}
	filter.MarketCap.setQuery(&q, "marketCap")
	filter.FullyDilutedValuation.setQuery(&q, "fullyDilutedValuation")
	filter.Volume.setQuery(&q, "volume")
	filter.PriceChange1h.setQuery(&q, "priceChange1h")
	filter.PriceChange1d.setQuery(&q, "priceChange1d")
	filter.PriceChange7d.setQuery(&q, "priceChange7d")
	filter.AvailableSupply.setQuery(&q, "availableSupply")
	filter.TotalSupply.setQuery(&q, "totalSupply")
	filter.Rank.setQuery(&q, "rank")
	filter.Price.setQuery(&q, "price")
	if filter.RiskScore.setQuery(&q, "riskScore") {
		t := true
		filter.IncludeRiskScore = &t
	}
	if filter.IncludeRiskScore != nil {
		q.Add("includeRiskScore", strconv.FormatBool(*filter.IncludeRiskScore))
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/coins?"+q.Encode(), nil)
	res, err := do[CoinResult](c, req)
	return res, err
}

type CoinChartResult struct {
	CoinID       string      `json:"coinId"`
	Chart        [][]float64 `json:"chart"`
	ErrorMessage string      `json:"errorMessage"`
}

// CoinCharts return historical chart data for multiple cryptocurrencies.
func (c *Client) CoinCharts(ctx context.Context, period ChartType, coinIDs []string) ([]CoinChartResult, error) {
	q := url.Values{}
	q.Add("period", string(period))
	q.Add("coinIds", strings.Join(coinIDs, ","))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/coins/charts?"+q.Encode(), nil)
	res, err := do[[]CoinChartResult](c, req)
	return res, err
}

// CoinByID return detailed information about a specific cryptocurrency using coinId.
func (c *Client) CoinByID(ctx context.Context, coinID, currency string) (Coin, error) {
	q := url.Values{}
	if currency != "" {
		q.Add("currency", currency)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		host+"/coins/"+url.QueryEscape(coinID)+"?"+q.Encode(), nil)
	res, err := do[Coin](c, req)
	return res, err
}

// CoinChartByID return historical chart data for a specific cryptocurrency using coinId.
func (c *Client) CoinChartByID(ctx context.Context, period ChartType, coinID string) ([][]float64, error) {
	q := url.Values{}
	q.Add("period", string(period))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		host+"/coins/"+url.QueryEscape(coinID)+"/charts?"+q.Encode(), nil)
	res, err := do[[][]float64](c, req)
	return res, err
}

// CoinAvgPrice return the historical average price of a specific cryptocurrency for a given date.
func (c *Client) CoinAvgPrice(ctx context.Context, coinID string, timestamp int64) (TopCurrency, error) {
	q := url.Values{}
	q.Add("coinId", coinID)
	q.Add("timestamp", strconv.FormatInt(timestamp, 10))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/coins/price/avg?"+q.Encode(), nil)
	res, err := do[TopCurrency](c, req)
	return res, err
}

// CoinExchangePrice return historical price data for a specific cryptocurrency on a selected exchange.
func (c *Client) CoinExchangePrice(ctx context.Context, exchange, fromCurrency, toCurrency string, timestamp int64) (float64, error) {
	q := url.Values{}
	q.Add("exchange", exchange)
	q.Add("from", fromCurrency)
	q.Add("to", toCurrency)
	q.Add("timestamp", strconv.FormatInt(timestamp, 10))
	type result struct {
		Price float64 `json:"price"`
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/coins/price/exchange?"+q.Encode(), nil)
	res, err := do[result](c, req)
	return res.Price, err
}

// TickerExchanges return a comprehensive list of cryptocurrency exchanges supported by CoinStats.
func (c *Client) TickerExchanges(ctx context.Context) ([]Exchange, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/tickers/exchanges", nil)
	res, err := do[[]Exchange](c, req)
	return res, err
}

type TickerMarketsFilter struct {
	Page         int
	Limit        int
	Exchange     string
	FromCoinID   string
	ToCoinID     string
	CoinID       string
	OnlyVerified *bool
}

type TickerMarketResult struct {
	Meta   Meta `json:"meta"`
	Result []struct {
		CreatedAt  time.Time `json:"_created_at"`
		UpdatedAt  time.Time `json:"_updated_at"`
		Exchange   string    `json:"exchange"`
		From       string    `json:"from"`
		To         string    `json:"to"`
		Pair       string    `json:"pair"`
		Price      float64   `json:"price"`
		PairPrice  float64   `json:"pairPrice"`
		Volume     float64   `json:"volume"`
		PairVolume float64   `json:"pairVolume"`
	} `json:"result"`
}

// TickerMarkets return a list of tickers for a specific cryptocurrency across multiple exchanges.
func (c *Client) TickerMarkets(ctx context.Context, filter TickerMarketsFilter) (TickerMarketResult, error) {
	q := url.Values{}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	if filter.Exchange != "" {
		q.Add("exchange", filter.Exchange)
	}
	if filter.FromCoinID != "" {
		q.Add("fromCoinId", filter.FromCoinID)
	}
	if filter.ToCoinID != "" {
		q.Add("toCoinId", filter.ToCoinID)
	}
	if filter.CoinID != "" {
		q.Add("coinId", filter.CoinID)
	}
	if filter.OnlyVerified != nil {
		q.Add("onlyVerified", strconv.FormatBool(*filter.OnlyVerified))
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/tickers/markets?"+q.Encode(), nil)
	res, err := do[TickerMarketResult](c, req)
	return res, err
}

// Fiats return detailed information about fiat currencies.
func (c *Client) Fiats(ctx context.Context) ([]Fiat, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/fiats", nil)
	res, err := do[[]Fiat](c, req)
	return res, err
}

// Markets return current global cryptocurrency market data.
func (c *Client) Markets(ctx context.Context) (Market, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/markets", nil)
	res, err := do[Market](c, req)
	return res, err
}

// Currencies return the complete list of supported fiat currencies.
func (c *Client) Currencies(ctx context.Context) (map[string]float64, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/currencies", nil)
	type result struct {
		Result map[string]float64 `json:"result"`
	}
	res, err := do[result](c, req)
	return res.Result, err
}
