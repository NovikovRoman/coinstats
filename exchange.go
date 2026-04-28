package coinstats

import (
	"bytes"
	"context"
	"encoding/json"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ExchangeResult struct {
	ConnectionID     string `json:"connectionId"`
	Name             string `json:"name"`
	Icon             string `json:"icon"`
	ConnectionFields []struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	} `json:"connectionFields"`
}

// Exchanges return the list of exchange portfolio connections supported by CoinStats.
func (c *Client) Exchanges(ctx context.Context) ([]ExchangeResult, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/exchange/support", nil)
	res, err := do[[]ExchangeResult](c, req)
	return res, err
}

type ExchangeBalanceResult struct {
	Balances []struct {
		CoinID   string  `json:"coinId"`
		Amount   float64 `json:"amount"`
		Price    float64 `json:"price"`
		PriceBtc float64 `json:"priceBtc"`
	} `json:"balances"`
	Portfolio struct {
		ID           string `json:"id"`
		Status       string `json:"status"`
		ConnectionID string `json:"connectionId"`
	} `json:"portfolio"`
}

// ExchangeBalance return cryptocurrency exchange balances.
func (c *Client) ExchangeBalance(ctx context.Context, connectionID, key, secret string, connFields map[string]string) (ExchangeBalanceResult, error) {
	fields := make(map[string]string, len(connFields)+2)
	maps.Copy(fields, connFields)
	fields["apiKey"] = key
	fields["apiSecret"] = secret
	data := struct {
		ConnectionID     string            `json:"connectionId"`
		ConnectionFields map[string]string `json:"connectionFields"`
	}{
		ConnectionID:     connectionID,
		ConnectionFields: fields,
	}
	b, _ := json.Marshal(data)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, host+"/exchange/balance", bytes.NewReader(b))
	res, err := do[ExchangeBalanceResult](c, req)
	return res, err
}

// ExchangeStatus return syncing status of the exchange portfolio, indicating whether the portfolio is fully synced with the exchange or still in progress.
func (c *Client) ExchangeStatus(ctx context.Context, portfolioID string) (SyncStatus, error) {
	q := url.Values{}
	q.Add("portfolioId", portfolioID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/exchange/status?"+q.Encode(), nil)
	type result struct {
		Status SyncStatus `json:"status"`
	}
	res, err := do[result](c, req)
	return res.Status, err
}

type ExchangeTransactionFilter struct {
	Page     int
	Limit    int
	DateFrom time.Time
	DateTo   time.Time
	Currency string
	Types    []string
}

type ExchangeTransactionResult struct {
	Meta   MetaShort     `json:"meta"`
	Result []Transaction `json:"result"`
}

// ExchangeTransactions return transaction data for a specific exchange by portfolioId.
func (c *Client) ExchangeTransactions(ctx context.Context, portfolioID string, filter ExchangeTransactionFilter) (ExchangeTransactionResult, error) {
	q := url.Values{}
	q.Add("portfolioId", portfolioID)
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	if !filter.DateFrom.IsZero() {
		q.Add("from", filter.DateFrom.Format(time.RFC3339Nano))
	}
	if !filter.DateTo.IsZero() {
		q.Add("to", filter.DateTo.Format(time.RFC3339Nano))
	}
	if filter.Currency != "" {
		q.Add("currency", filter.Currency)
	}
	if len(filter.Types) > 0 {
		q.Add("types", strings.Join(filter.Types, ","))
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/exchange/transactions?"+q.Encode(), nil)
	res, err := do[ExchangeTransactionResult](c, req)
	return res, err
}

// ExchangeChart return exchange chart data for specific time ranges displayed on the CoinStats website.
func (c *Client) ExchangeChart(ctx context.Context, portfolioID string, t ChartType) ([][]float64, error) {
	q := url.Values{}
	q.Add("portfolioId", portfolioID)
	q.Add("type", string(t))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/exchange/chart?"+q.Encode(), nil)
	res, err := do[[][]float64](c, req)
	return res, err
}

// ExchangeSync initiate syncing process for the given exchange portfolio by portfolioId.
func (c *Client) ExchangeSync(ctx context.Context, portfolioID string) (bool, error) {
	q := url.Values{}
	q.Add("portfolioId", portfolioID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPatch, host+"/exchange/sync?"+q.Encode(), nil)
	type result struct {
		Success bool `json:"success"`
	}
	res, err := do[result](c, req)
	return res.Success, err
}

type ExchangePLFilter struct {
	Page   int
	Limit  int
	CoinID string
}

type ExchangePLResult struct {
	Result []struct {
		Count           float64     `json:"count"`
		Coin            Сoin        `json:"coin"`
		Price           TopCurrency `json:"price"`
		ProfitPercent   Profit      `json:"profitPercent"`
		Profit          Profit      `json:"profit"`
		AverageBuy      Profit      `json:"averageBuy"`
		AverageSell     Profit      `json:"averageSell"`
		LiquidityScore  float64     `json:"liquidityScore"`
		VolatilityScore float64     `json:"volatilityScore"`
		MarketCapScore  float64     `json:"marketCapScore"`
		RiskScore       float64     `json:"riskScore"`
		AvgChange       float64     `json:"avgChange"`
		TotalCost       TopCurrency `json:"totalCost"`
	} `json:"result"`
	Summary PLSummary `json:"summary"`
}

// ExchangePL return profit/loss data for a specific exchange portfolio.
func (c *Client) ExchangePL(ctx context.Context, portfolioID string, filter ExchangePLFilter) (ExchangePLResult, error) {
	q := url.Values{}
	q.Add("portfolioId", portfolioID)
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	if filter.CoinID != "" {
		q.Add("coinId", filter.CoinID)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/exchange/pl?"+q.Encode(), nil)
	res, err := do[ExchangePLResult](c, req)
	return res, err
}
