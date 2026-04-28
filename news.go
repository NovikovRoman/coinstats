package coinstats

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type NewsSourceResult struct {
	Sourcename string    `json:"sourcename"`
	WebURL     string    `json:"weburl"`
	FeedURL    string    `json:"feedurl"`
	CoinID     string    `json:"coinid"`
	Logo       string    `json:"logo"`
	SourceImg  string    `json:"sourceImg"`
	CreatedAt  time.Time `json:"_created_at"`
	UpdatedAt  time.Time `json:"_updated_at"`
}

// NewsSources return the list of news sources.
func (c *Client) NewsSources(ctx context.Context) ([]NewsSourceResult, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/news/sources", nil)
	res, err := do[[]NewsSourceResult](c, req)
	return res, err
}

type NewsFilter struct {
	Page     int
	Limit    int
	DateFrom time.Time
	DateTo   time.Time
}

type NewsResult struct {
	ID             string         `json:"id"`
	FeedDate       int64          `json:"feedDate"`
	Source         string         `json:"source"`
	Title          string         `json:"title"`
	IsFeatured     bool           `json:"isFeatured"`
	Link           string         `json:"link"`
	SourceLink     string         `json:"sourceLink"`
	ImgURL         string         `json:"imgUrl"`
	ReactionsCount map[string]int `json:"reactionsCount"`
	Reactions      []string       `json:"reactions"`
	ShareURL       string         `json:"shareURL"`
	RelatedCoins   []string       `json:"relatedCoins"`
	Content        bool           `json:"content"`
	BigImg         bool           `json:"bigImg"`
	SearchKeyWords []string       `json:"searchKeyWords"`
	Description    string         `json:"description"`
	Coins          []struct {
		CoinKeyWords      string  `json:"coinKeyWords"`
		CoinPercent       float64 `json:"coinPercent"`
		CoinTitleKeyWords string  `json:"coinTitleKeyWords"`
		CoinNameKeyWords  string  `json:"coinNameKeyWords"`
		CoinIdKeyWords    string  `json:"coinIdKeyWords"`
	} `json:"coins"`
}

// News return the list of cryptocurrency news articles with pagination.
func (c *Client) News(ctx context.Context, filter NewsFilter) ([]NewsResult, error) {
	q := url.Values{}
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
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/news?"+q.Encode(), nil)
	type result struct {
		Result []NewsResult `json:"result"`
	}
	res, err := do[result](c, req)
	return res.Result, err
}

type NewsType string

const (
	NewsTypeHandpicked NewsType = "handpicked"
	NewsTypeTrending   NewsType = "trending"
	NewsTypeLatest     NewsType = "latest"
	NewsTypeBullish    NewsType = "bullish"
	NewsTypeBearish    NewsType = "bearish"
)

type NewsByTypeFilter struct {
	Page  int
	Limit int
}

// NewsByType return cryptocurrency news articles based on a specific type.
func (c *Client) NewsByType(ctx context.Context, t NewsType, filter NewsByTypeFilter) ([]NewsResult, error) {
	q := url.Values{}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		host+"/news/type/"+url.QueryEscape(string(t))+"?"+q.Encode(), nil)
	res, err := do[[]NewsResult](c, req)
	return res, err
}

// NewsByID return a news article by id.
func (c *Client) NewsByID(ctx context.Context, id string) (NewsResult, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/news/"+url.QueryEscape(string(id)), nil)
	res, err := do[NewsResult](c, req)
	return res, err
}
