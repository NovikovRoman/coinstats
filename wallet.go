package coinstats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const ConnectionIDAll = "all"

// Blockchains return the list of blockchains supported by CoinStats.
func (c *Client) Blockchains(ctx context.Context) ([]Blockchain, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/wallet/blockchains", nil)
	res, err := do[[]Blockchain](c, req)
	return res, err
}

// WalletBalance return cryptocurrency balances for any blockchain wallet.
func (c *Client) WalletBalance(ctx context.Context, w Wallet) ([]WalletBalance, error) {
	if w.ConnectionID == "" && w.Blockchain == "" {
		w.ConnectionID = ConnectionIDAll
	}
	q := url.Values{
		"address":      []string{w.Address},
		"connectionId": []string{w.ConnectionID},
		"blockchain":   []string{w.Blockchain},
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/wallet/balance?"+q.Encode(), nil)
	res, err := do[[]WalletBalance](c, req)
	return res, err
}

type WalletBalancesFilter struct {
	Address     string
	Blockchains []string
	Wallets     []WalletShort
}

// WalletBalances return cryptocurrency balances for multiple wallets.
func (c *Client) WalletBalances(ctx context.Context, filter WalletBalancesFilter) ([]WalletMultiBalance, error) {
	blockchains := map[string]struct{}{}
	for _, b := range filter.Blockchains {
		blockchains[b] = struct{}{}
	}

	wallets := []string{}
	for _, w := range filter.Wallets {
		if w.Address == "" {
			return []WalletMultiBalance{}, ErrBadWalletData
		}
		if w.ConnectionID == "" {
			w.ConnectionID = ConnectionIDAll
		}
		wallets = append(wallets, fmt.Sprintf("%s:%s", w.ConnectionID, w.Address))
	}

	q := url.Values{}
	q.Add("address", filter.Address)
	filter.Blockchains = make([]string, 0, len(blockchains))
	for b := range blockchains {
		filter.Blockchains = append(filter.Blockchains, b)
	}
	q.Add("blockchain", strings.Join(filter.Blockchains, ","))
	q.Add("wallets", strings.Join(wallets, ","))

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/wallet/balances?"+q.Encode(), nil)
	res, err := do[[]WalletMultiBalance](c, req)
	return res, err
}

// WalletSyncStatus return the syncing status of the provided wallet address with the blockchain network.
func (c *Client) WalletSyncStatus(ctx context.Context, w Wallet) (SyncStatus, error) {
	q := url.Values{}
	q.Add("address", w.Address)
	q.Add("connectionId", w.ConnectionID)
	q.Add("blockchain", w.Blockchain)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/wallet/status?"+q.Encode(), nil)
	type result struct {
		Status SyncStatus `json:"status"`
	}
	res, err := do[result](c, req)
	return res.Status, err
}

type WalletTransactionFilter struct {
	Wallet   Wallet
	Page     int
	Limit    int
	DateFrom time.Time
	DateTo   time.Time
	Currency string
	Types    []string
	TxID     string
	CoinID   string
	Wallets  []WalletShort
}

type WalletTransactionsResult struct {
	Meta   MetaShort     `json:"meta"`
	Result []Transaction `json:"result"`
}

// WalletTransactions return transaction data for wallet addresses.
func (c *Client) WalletTransactions(ctx context.Context, filter WalletTransactionFilter) (WalletTransactionsResult, error) {
	wallets := []string{}
	for _, w := range filter.Wallets {
		if w.Address == "" {
			return WalletTransactionsResult{}, ErrBadWalletData
		}
		if w.ConnectionID == "" {
			w.ConnectionID = ConnectionIDAll
		}
		wallets = append(wallets, fmt.Sprintf("%s:%s", w.ConnectionID, w.Address))
	}

	q := url.Values{}
	if !filter.DateFrom.IsZero() {
		q.Add("from", filter.DateFrom.Format(time.RFC3339Nano))
	}
	if !filter.DateTo.IsZero() {
		q.Add("to", filter.DateTo.Format(time.RFC3339Nano))
	}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	if filter.Currency != "" {
		q.Add("currency", filter.Currency)
	}
	if filter.Wallet.Address != "" {
		q.Add("address", filter.Wallet.Address)
		q.Add("connectionId", filter.Wallet.ConnectionID)
		q.Add("blockchain", filter.Wallet.Blockchain)
	}
	q.Add("wallets", strings.Join(wallets, ","))
	if len(filter.Types) > 0 {
		q.Add("types", strings.Join(filter.Types, ","))
	}
	if filter.CoinID != "" {
		q.Add("coinId", filter.CoinID)
	}
	if filter.TxID != "" {
		q.Add("txId", filter.TxID)
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/wallet/transactions?"+q.Encode(), nil)
	res, err := do[WalletTransactionsResult](c, req)
	return res, err
}

// WalletTransactionsSync initiate syncing process to update transaction data.
func (c *Client) WalletTransactionsSync(ctx context.Context, wallets []Wallet) (SyncStatus, error) {
	var body io.Reader
	data := struct {
		Wallets []Wallet `json:"wallets"`
	}{
		Wallets: wallets,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return UnknownStatus, err
	}
	body = bytes.NewReader(b)

	type result struct {
		Status SyncStatus `json:"status"`
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPatch, host+"/wallet/transactions", body)
	res, err := do[result](c, req)
	return res.Status, err
}

// WalletChart return wallet chart data for specific time ranges.
func (c *Client) WalletChart(ctx context.Context, t ChartType, w Wallet) ([][]float64, error) {
	q := url.Values{}
	q.Add("address", w.Address)
	q.Add("connectionId", w.ConnectionID)
	q.Add("blockchain", w.Blockchain)
	q.Add("type", string(t))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/wallet/chart?"+q.Encode(), nil)
	type result struct {
		Result [][]float64 `json:"result"`
	}
	res, err := do[result](c, req)
	return res.Result, err
}

type WalletChartsResult struct {
	Data          [][]float64 `json:"data"`
	WalletAddress string      `json:"walletAddress"`
	ConnectionID  string      `json:"connectionId"`
	Blockchain    string      `json:"blockchain"`
	Message       string      `json:"message"`
}

// WalletCharts return chart data for multiple wallet addresses across various networks.
func (c *Client) WalletCharts(ctx context.Context, t ChartType, w []WalletShort, aggregated bool) ([]WalletChartsResult, error) {
	q := url.Values{}
	q.Add("type", string(t))
	wallets := make([]string, 0, len(w))
	for _, ww := range w {
		connID := ww.ConnectionID
		if connID == "" {
			connID = ConnectionIDAll
		}
		wallets = append(wallets, fmt.Sprintf("%s:%s", connID, ww.Address))
	}
	q.Add("wallets", strings.Join(wallets, ","))
	q.Add("aggregated", strconv.FormatBool(aggregated))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/wallet/charts?"+q.Encode(), nil)
	res, err := do[[]WalletChartsResult](c, req)
	return res, err
}

// WalletDefi return comprehensive DeFi portfolio data, including staking, liquidity pool (LP), and yield farming activities.
func (c *Client) WalletDefi(ctx context.Context, w Wallet) (Defi, error) {
	q := url.Values{}
	q.Add("address", w.Address)
	q.Add("connectionId", w.ConnectionID)
	q.Add("blockchain", w.Blockchain)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/wallet/defi?"+q.Encode(), nil)
	res, err := do[Defi](c, req)
	return res, err
}

type WalletPLItem struct {
	Count         float64     `json:"count"`
	Coin          Сoin        `json:"coin"`
	AverageBuy    Profit      `json:"averageBuy"`
	AverageSell   Profit      `json:"averageSell"`
	Price         TopCurrency `json:"price"`
	Profit        Profit      `json:"profit"`
	ProfitPercent Profit      `json:"profitPercent"`
	TotalCost     TopCurrency `json:"totalCost"`
}

type WalletPLResult struct {
	Result  []WalletPLItem `json:"result"`
	Summary PLSummary      `json:"summary"`
}

type WalletPLFilter struct {
	Wallet Wallet
	CoinID string
	Page   int
	Limit  int
}

// WalletPL return profit/loss data for a wallet.
func (c *Client) WalletPL(ctx context.Context, filter WalletPLFilter) (WalletPLResult, error) {
	q := url.Values{}
	if filter.Wallet.Address != "" {
		q.Add("address", filter.Wallet.Address)
		q.Add("connectionId", filter.Wallet.ConnectionID)
		q.Add("blockchain", filter.Wallet.Blockchain)
	}
	if filter.CoinID != "" {
		q.Add("coinId", filter.CoinID)
	}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/wallet/pl?"+q.Encode(), nil)
	res, err := do[WalletPLResult](c, req)
	return res, err
}
