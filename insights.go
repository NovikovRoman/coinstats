package coinstats

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// BTCDominance return Bitcoin market dominance data showing BTC’s percentage share of the total cryptocurrency market capitalization over a specified time period.
func (c *Client) BTCDominance(ctx context.Context, t ChartType) ([][]float64, error) {
	q := url.Values{}
	q.Add("type", string(t))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/insights/btc-dominance?"+q.Encode(), nil)
	type result struct {
		Data [][]float64 `json:"data"`
	}
	res, err := do[result](c, req)
	return res.Data, err
}

type fearAndGreedValue struct {
	Value               int    `json:"value"`
	ValueClassification string `json:"value_classification"`
	Timestamp           int64  `json:"timestamp"`
}

type FearAndGreedResult struct {
	Name string `json:"name"`
	Now  struct {
		Value               int       `json:"value"`
		ValueClassification string    `json:"value_classification"`
		Timestamp           int64     `json:"timestamp"`
		UpdateTime          time.Time `json:"update_time"`
	} `json:"now"`
	Yesterday fearAndGreedValue `json:"yesterday"`
	LastWeek  fearAndGreedValue `json:"lastWeek"`
}

// FearAndGreed return the Crypto Fear & Greed Index.
func (c *Client) FearAndGreed(ctx context.Context) (FearAndGreedResult, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/insights/fear-and-greed", nil)
	res, err := do[FearAndGreedResult](c, req)
	return res, err
}

type FearAndGreedChartResult struct {
	Name string `json:"name"`
	Data []struct {
		Value               int    `json:"value"`
		ValueClassification string `json:"value_classification"`
		Timestamp           string `json:"timestamp"`
	} `json:"data"`
}

// FearAndGreedChart return historical data for the Crypto Fear & Greed Index.
func (c *Client) FearAndGreedChart(ctx context.Context) (FearAndGreedChartResult, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/insights/fear-and-greed/chart", nil)
	res, err := do[FearAndGreedChartResult](c, req)
	return res, err
}

type RainbowChartResult struct {
	Price   float64   `json:"price"`
	TimeRaw string    `json:"time"`
	Time    time.Time `json:"-"`
}

// RainbowChart return Rainbow Chart data.
func (c *Client) RainbowChart(ctx context.Context, coinID string) ([]RainbowChartResult, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		host+"/insights/rainbow-chart/"+url.QueryEscape(coinID), nil)
	res, err := do[[]RainbowChartResult](c, req)
	for i := range res {
		if res[i].TimeRaw != "" {
			res[i].Time, _ = time.Parse(time.DateOnly, res[i].TimeRaw)
		}
	}
	return res, err
}
