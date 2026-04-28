package coinstats

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type PortfolioResult struct {
	ID   string `json:"portfolioId"`
	Name string `json:"portfolioName"`
}

// PortfolioList return list of all API-connected portfolios
func (c *Client) PortfolioList(ctx context.Context) ([]PortfolioResult, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/portfolio/list", nil)
	res, err := do[[]PortfolioResult](c, req)
	return res, err
}

type PortfolioConnectResult struct {
	PortfolioID  string                   `json:"portfolioId"`
	ConnectionID string                   `json:"connectionId"`
	Status       PortfolioConnectedStatus `json:"status"`
}

// PortfolioWallet connect a wallet to your account and create a tracked portfolio
func (c *Client) PortfolioWallet(ctx context.Context, address, connectionID, name string) (PortfolioConnectResult, error) {
	data := struct {
		Address      string `json:"address"`
		ConnectionID string `json:"connectionId"`
		Name         string `json:"name"`
	}{
		Address:      address,
		ConnectionID: connectionID,
		Name:         name,
	}
	b, _ := json.Marshal(data)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, host+"/portfolio/wallet", bytes.NewReader(b))
	res, err := do[PortfolioConnectResult](c, req)
	return res, err
}

// PortfolioExchange connect an exchange account to your account and create a tracked portfolio
func (c *Client) PortfolioExchange(ctx context.Context, apiKey, apiSecret, connectionID, name string) (PortfolioConnectResult, error) {
	type fields struct {
		ApiKey    string `json:"apiKey"`
		ApiSecret string `json:"apiSecret"`
	}
	type data struct {
		ConnectionID     string `json:"connectionId"`
		Name             string `json:"name"`
		ConnectionFields fields `json:"connectionFields"`
	}
	d := data{
		ConnectionID: connectionID,
		Name:         name,
		ConnectionFields: fields{
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
		},
	}
	b, _ := json.Marshal(d)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, host+"/portfolio/exchange", bytes.NewReader(b))
	res, err := do[PortfolioConnectResult](c, req)
	return res, err
}

type PortfolioValueResult struct {
	UnrealizedProfitLossPercent float64 `json:"unrealizedProfitLossPercent"`
	RealizedProfitLossPercent   float64 `json:"realizedProfitLossPercent"`
	AllTimeProfitLossPercent    float64 `json:"allTimeProfitLossPercent"`
	TotalValue                  float64 `json:"totalValue"`
	DefiValue                   float64 `json:"defiValue"`
	TotalCost                   float64 `json:"totalCost"`
	UnrealizedProfitLoss        float64 `json:"unrealizedProfitLoss"`
	RealizedProfitLoss          float64 `json:"realizedProfitLoss"`
	AllTimeProfitLoss           float64 `json:"allTimeProfitLoss"`
}

// PortfolioValue return detailed information about portfolio profit/loss.
func (c *Client) PortfolioValue(ctx context.Context, passcode, portfolioID, currency string) (PortfolioValueResult, error) {
	q := url.Values{}
	q.Add("currency", currency)
	if !c.hasShareToken() {
		q.Add("portfolioId", portfolioID)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/portfolio/value?"+q.Encode(), nil)
	c.addShareToken(req)
	req.Header.Set("passcode", passcode)
	res, err := do[PortfolioValueResult](c, req)
	return res, err
}

type PortfolioCoinsFilter struct {
	Page             int
	Limit            int
	IncludeRiskScore *bool
}

type PortfolioCoinResult struct {
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
}

// PortfolioCoins return detailed information about portfolio profit/loss.
func (c *Client) PortfolioCoins(ctx context.Context, passcode, portfolioID string, filter PortfolioCoinsFilter) (PortfolioCoinResult, error) {
	q := url.Values{}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	if filter.IncludeRiskScore != nil {
		q.Add("includeRiskScore", strconv.FormatBool(*filter.IncludeRiskScore))
	}
	if !c.hasShareToken() {
		q.Add("portfolioId", portfolioID)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/portfolio/coins?"+q.Encode(), nil)
	c.addShareToken(req)
	req.Header.Set("passcode", passcode)
	res, err := do[PortfolioCoinResult](c, req)
	return res, err
}

// PortfolioChart return historical performance data to visualize your portfolio’s growth over time.
func (c *Client) PortfolioChart(ctx context.Context, passcode, portfolioID string, t ChartType) ([][]float64, error) {
	q := url.Values{}
	q.Add("type", string(t))
	if !c.hasShareToken() {
		q.Add("portfolioId", portfolioID)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/portfolio/chart?"+q.Encode(), nil)
	c.addShareToken(req)
	req.Header.Set("passcode", passcode)
	type result struct {
		Result [][]float64 `json:"result"`
	}
	res, err := do[result](c, req)
	return res.Result, err
}

type PortfolioTransactionFilter struct {
	Page     int
	Limit    int
	Currency string
	CoinID   string
}

type PortfolioTransactionResult struct {
	Meta MetaShort              `json:"meta"`
	Data []PortfolioTransaction `json:"data"`
}

// PortfolioTransactions return a detailed history of all transactions in your portfolio.
func (c *Client) PortfolioTransactions(ctx context.Context, passcode, portfolioID string, filter PortfolioTransactionFilter) (PortfolioTransactionResult, error) {
	q := url.Values{}
	q.Add("currency", filter.Currency)
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	if filter.CoinID != "" {
		q.Add("coinId", filter.CoinID)
	}
	if !c.hasShareToken() {
		q.Add("portfolioId", portfolioID)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/portfolio/transactions?"+q.Encode(), nil)
	c.addShareToken(req)
	req.Header.Set("passcode", passcode)
	res, err := do[PortfolioTransactionResult](c, req)
	return res, err
}

type TransactionData struct {
	CoinID      string  `json:"coinId"`
	Count       float64 `json:"count"`
	Date        int64   `json:"date"`
	Price       float64 `json:"price"`
	PortfolioID string  `json:"portfolioId"`
	Currency    string  `json:"currency"`
	Notes       string  `json:"notes"`
}

// AddPortfolioTransaction add a new transaction to your manual portfolio.
func (c *Client) AddPortfolioTransaction(ctx context.Context, passcode, portfolioID string, t TransactionData) (PortfolioTransaction, error) {
	q := url.Values{}
	if !c.hasShareToken() {
		q.Add("portfolioId", portfolioID)
	}
	b, _ := json.Marshal(t)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, host+"/portfolio/transaction?"+q.Encode(), bytes.NewReader(b))
	c.addShareToken(req)
	req.Header.Set("passcode", passcode)
	res, err := do[PortfolioTransaction](c, req)
	return res, err
}

// PortfolioDefi return comprehensive DeFi portfolio data, including staking, liquidity pool (LP), and yield farming activities.
func (c *Client) PortfolioDefi(ctx context.Context, passcode, portfolioID string) (Defi, error) {
	q := url.Values{}
	if !c.hasShareToken() {
		q.Add("portfolioId", portfolioID)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/portfolio/defi?"+q.Encode(), nil)
	c.addShareToken(req)
	req.Header.Set("passcode", passcode)
	res, err := do[Defi](c, req)
	return res, err
}

type PortfolioSnapshotItemFilter struct {
	Page     int
	Limit    int
	DateFrom time.Time
	DateTo   time.Time
	CoinID   string
}

type PortfolioSnapshotItemResult struct {
	Result []struct {
		Date         time.Time `json:"date"`
		CoinBalances []struct {
			Symbol             string  `json:"symbol"`
			Icon               string  `json:"icon"`
			Quantity           float64 `json:"quantity"`
			Balance            float64 `json:"balance"`
			QuantityChange     float64 `json:"quantityChange"`
			BalanceChange      float64 `json:"balanceChange"`
			PricePerUnit       float64 `json:"pricePerUnit"`
			PricePercentChange float64 `json:"pricePercentChange"`
		} `json:"coinBalances"`
		TotalBalance              float64 `json:"totalBalance"`
		TotalBalanceChange        float64 `json:"totalBalanceChange"`
		TotalBalancePercentChange float64 `json:"totalBalancePercentChange"`
	} `json:"result"`
}

// PortfolioSnapshotItems return historical portfolio snapshot data with normalized coin balances and portfolio metrics.
func (c *Client) PortfolioSnapshotItems(ctx context.Context, passcode, portfolioID string, filter PortfolioSnapshotItemFilter) (PortfolioSnapshotItemResult, error) {
	q := url.Values{}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	if filter.CoinID != "" {
		q.Add("coinId", filter.CoinID)
	}
	if !filter.DateFrom.IsZero() {
		q.Add("from", filter.DateFrom.Format(time.DateOnly))
	}
	if !filter.DateTo.IsZero() {
		q.Add("to", filter.DateTo.Format(time.DateOnly))
	}
	if !c.hasShareToken() {
		q.Add("portfolioId", portfolioID)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/portfolio/snapshot/items?"+q.Encode(), nil)
	c.addShareToken(req)
	req.Header.Set("passcode", passcode)
	res, err := do[PortfolioSnapshotItemResult](c, req)
	return res, err
}

// PortfolioStatus return portfolio sync status
func (c *Client) PortfolioStatus(ctx context.Context, portfolioID string) (SyncStatus, error) {
	q := url.Values{
		"portfolioId": []string{portfolioID},
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/portfolio/status?"+q.Encode(), nil)
	type result struct {
		Status SyncStatus `json:"status"`
	}
	res, err := do[result](c, req)
	return res.Status, err
}

// PortfolioSync sync portfolio
func (c *Client) PortfolioSync(ctx context.Context, portfolioID string) (bool, error) {
	q := url.Values{
		"portfolioId": []string{portfolioID},
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPatch, host+"/portfolio/sync?"+q.Encode(), nil)
	type result struct {
		Success bool `json:"success"`
	}
	res, err := do[result](c, req)
	return res.Success, err
}

// PortfolioDelete delete portfolio
func (c *Client) PortfolioDelete(ctx context.Context, portfolioID string) error {
	q := url.Values{
		"portfolioId": []string{portfolioID},
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, host+"/portfolio?"+q.Encode(), nil)
	type result struct {
		Success bool `json:"success"`
	}
	_, err := do[result](c, req) //! в документации нет информации об успешном ответе
	return err
}
