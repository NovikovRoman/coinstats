package coinstats

import (
	"fmt"
	"net/url"
	"time"
)

type SyncStatus string

const (
	UnknownStatus SyncStatus = ""
	SyncingStatus SyncStatus = "syncing"
	SyncedStatus  SyncStatus = "synced"
)

type PortfolioConnectedStatus string

const (
	UnknownConnectedStatus PortfolioConnectedStatus = ""
	ConnectedStatus        PortfolioConnectedStatus = "connected"
	ExistingStatus         PortfolioConnectedStatus = "existing"
)

type ChartType string

const (
	ChartType24h ChartType = "24h"
	ChartType1w  ChartType = "1w"
	ChartType1m  ChartType = "1m"
	ChartType3m  ChartType = "3m"
	ChartType6m  ChartType = "6m"
	ChartType1y  ChartType = "1y"
	ChartTypeAll ChartType = "all"
)

type SortDir string

const (
	Asc  SortDir = "asc"
	Desc SortDir = "desc"
)

type CoinBase struct {
	ID                    string  `json:"id"`
	Icon                  string  `json:"icon"`
	Name                  string  `json:"name"`
	Symbol                string  `json:"symbol"`
	Rank                  int     `json:"rank"`
	Price                 float64 `json:"price"`
	PriceBtc              float64 `json:"priceBtc"`
	Volume                float64 `json:"volume"`
	MarketCap             float64 `json:"marketCap"`
	AvailableSupply       float64 `json:"availableSupply"`
	TotalSupply           float64 `json:"totalSupply"`
	FullyDilutedValuation float64 `json:"fullyDilutedValuation"`
	PriceChange1h         float64 `json:"priceChange1h"`
	PriceChange1d         float64 `json:"priceChange1d"`
	PriceChange1w         float64 `json:"priceChange1w"`
	WebsiteURL            string  `json:"websiteUrl"`
	RedditURL             string  `json:"redditUrl"`
	TwitterURL            string  `json:"twitterUrl"`
	ContractAddress       string  `json:"contractAddress"`
	ContractAddresses     []struct {
		Blockchain      string `json:"blockchain"`
		ContractAddress string `json:"contractAddress"`
	} `json:"contractAddresses"`
	Decimals        int      `json:"decimals"`
	Explorers       []string `json:"explorers"`
	LiquidityScore  float64  `json:"liquidityScore"`
	VolatilityScore float64  `json:"volatilityScore"`
	MarketCapScore  float64  `json:"marketCapScore"`
	RiskScore       float64  `json:"riskScore"`
	AvgChange       float64  `json:"avgChange"`
}

type Coin struct {
	CoinBase
	Slug     string `json:"slug"`
	IsStable bool   `json:"isStable"`
	Color    string `json:"color"`
}

type Wallet struct {
	Address      string `json:"address"`
	ConnectionID string `json:"connectionId,omitempty"`
	Blockchain   string `json:"blockchain,omitempty"`
}

type WalletShort struct {
	Address      string `json:"address"`
	ConnectionID string `json:"connectionId"`
}

type WalletBalance struct {
	CoinID          string  `json:"coinId"`
	Amount          float64 `json:"amount"`
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	Price           float64 `json:"price"`
	PriceBTC        float64 `json:"priceBtc"`
	ImgURL          string  `json:"imgUrl"`
	PriceChange24h  float64 `json:"pCh24h"`
	Rank            int     `json:"rank"`
	Volume          float64 `json:"volume"`
	Chain           string  `json:"chain"`
	Decimals        int64   `json:"decimals"`
	ContractAddress string  `json:"contractAddress"`
}

type WalletMultiBalance struct {
	Blockchain   string          `json:"blockchain"`
	Address      string          `json:"address"`
	ConnectionID string          `json:"connectionId"`
	Balances     []WalletBalance `json:"balances"`
}

type Blockchain struct {
	Name         string `json:"name"`
	ConnectionID string `json:"connectionId"`
	Chain        string `json:"chain"`
	Icon         string `json:"icon"`
}

type Thresholds struct {
	GreaterThan float64
	Equals      float64
	LessThan    float64
}

func (t Thresholds) setQuery(q *url.Values, name string) bool {
	ok := false
	if t.GreaterThan != 0 {
		ok = true
		q.Add(name+"~greaterThan", fmt.Sprintf("%f", t.GreaterThan))
	}
	if t.Equals != 0 {
		ok = true
		q.Add(name+"~equals", fmt.Sprintf("%f", t.Equals))
	}
	if t.LessThan != 0 {
		ok = true
		q.Add(name+"~lessThan", fmt.Sprintf("%f", t.LessThan))
	}
	return ok
}

type Exchange struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Rank      int     `json:"rank"`
	Change24h float64 `json:"change24h"`
	Volume24h float64 `json:"volume24h"`
	Volume7d  float64 `json:"volume7d"`
	Volume1m  float64 `json:"volume1m"`
	Icon      string  `json:"icon"`
	URL       string  `json:"url"`
}

type Fiat struct {
	Name     string  `json:"name"`
	Rate     float64 `json:"rate"`
	Symbol   string  `json:"symbol"`
	ImageURL string  `json:"imageUrl"`
}

type Market struct {
	MarketCap          float64 `json:"marketCap"`
	Volume             float64 `json:"volume"`
	BtcDominance       float64 `json:"btcDominance"`
	MarketCapChange    float64 `json:"marketCapChange"`
	VolumeChange       float64 `json:"volumeChange"`
	BtcDominanceChange float64 `json:"btcDominanceChange"`
}

type PortfolioTransaction struct {
	TransactionType string    `json:"transactionType"`
	Date            time.Time `json:"date"`
	ProfitLoss      struct {
		Profit        float64 `json:"profit"`
		ProfitPercent float64 `json:"profitPercent"`
		CurrentValue  float64 `json:"currentValue"`
	} `json:"profitLoss"`
	Fee struct {
		Coin        coin    `json:"coin"`
		Count       float64 `json:"count"`
		TotalWorth  float64 `json:"totalWorth"`
		Price       float64 `json:"price"`
		ToAddress   string  `json:"toAddress"`
		FromAddress string  `json:"fromAddress"`
	} `json:"fee"`
	PortfolioInfo struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	} `json:"portfolioInfo"`
	CoinData struct {
		Identifier   string  `json:"identifier"`
		Count        float64 `json:"count"`
		Symbol       string  `json:"symbol"`
		TotalWorth   float64 `json:"totalWorth"`
		CurrentValue float64 `json:"currentValue"`
	} `json:"coinData"`
	Transfers []struct {
		TransferType string `json:"transferType"`
		Items        []struct {
			Coin        coin    `json:"coin"`
			Count       float64 `json:"count"`
			ToAddress   string  `json:"toAddress"`
			FromAddress string  `json:"fromAddress"`
			TotalWorth  float64 `json:"totalWorth"`
			Price       float64 `json:"price"`
		} `json:"items"`
	} `json:"transfers"`
	Note struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	} `json:"note"`
}

type Transaction struct {
	Type        string    `json:"type"`
	Date        time.Time `json:"date"` // "2025-06-07T11:58:11.000Z"
	MainContent struct {
		CoinIcons  []string `json:"coinIcons"`
		CoinAssets []string `json:"coinAssets"`
	} `json:"mainContent"`
	CoinData struct {
		Count        float64 `json:"count"`
		Symbol       string  `json:"symbol"`
		CurrentValue float64 `json:"currentValue"`
	} `json:"coinData"`
	ProfitLoss struct {
		Profit        float64 `json:"profit"`
		ProfitPercent float64 `json:"profitPercent"`
		CurrentValue  float64 `json:"currentValue"`
	} `json:"profitLoss"`
	Transactions []struct {
		Action string `json:"action"`
		Items  []struct {
			ID         string  `json:"id"`
			Count      float64 `json:"count"`
			TotalWorth float64 `json:"totalWorth"`
			Coin       struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Symbol string `json:"symbol"`
				Icon   string `json:"icon"`
			} `json:"coin"`
		} `json:"items"`
	} `json:"transactions"`
	Fee struct {
		Coin struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
			Icon   string `json:"icon"`
		} `json:"coin"`
		Count      float64 `json:"count"`
		TotalWorth float64 `json:"totalWorth"`
	} `json:"fee"`
	Hash struct {
		ID          string `json:"id"`
		ExplorerURL string `json:"explorerUrl"`
	} `json:"hash"`
}

type Defi struct {
	TotalAssets TopCurrency `json:"totalAssets"`
	Protocols   []struct {
		ID         string      `json:"id"`
		Name       string      `json:"name"`
		Logo       string      `json:"logo"`
		TotalValue TopCurrency `json:"totalValue"`
		Blockchain struct {
			Name string `json:"name"`
			Icon string `json:"icon"`
		} `json:"blockchain"`
		Chain       string `json:"chain"`
		Investments []struct {
			ID      string      `json:"id"`
			Name    string      `json:"name"`
			Value   TopCurrency `json:"value"`
			Symbols string      `json:"symbols"`
			Assets  []struct {
				Address string      `json:"address"`
				Title   string      `json:"title"`
				Amount  float64     `json:"amount"`
				Symbol  string      `json:"symbol"`
				Price   TopCurrency `json:"price"`
				CoinID  string      `json:"coinId"`
				Logo    string      `json:"logo"`
				Danger  bool        `json:"danger"`
			} `json:"assets"`
			PoolAddress    string    `json:"poolAddress"`
			HealthRate     int       `json:"healthRate"`
			HealthRateLink string    `json:"healthRateLink"`
			PoolProjectID  string    `json:"poolProjectId"`
			UnlockAt       time.Time `json:"unlockAt"`
			EndAt          time.Time `json:"endAt"`
			IsProxy        bool      `json:"isProxy"`
			DebtRatio      float64   `json:"debtRatio"`
			DebtRatioLink  string    `json:"debtRatioLink"`
		} `json:"investments"`
		WalletAddress string `json:"walletAddress"`
		URL           string `json:"url"`
	} `json:"protocols"`
}

type TopCurrency struct {
	USD float64 `json:"USD"`
	ETH float64 `json:"ETH"`
	BTC float64 `json:"BTC"`
}

type NftCollection struct {
	Address      string `json:"address"`
	Description  string `json:"description"`
	Name         string `json:"name"`
	BannerImg    string `json:"bannerImg"`
	Blockchain   string `json:"blockchain"`
	Img          string `json:"img"`
	RelevantURLs []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"relevantUrls"`
	Slug                   string    `json:"slug"`
	SlugCs                 string    `json:"slugCs"`
	Source                 string    `json:"source"`
	Verified               bool      `json:"verified"`
	Rank                   int       `json:"rank"`
	RankAll                int       `json:"rankAll"`
	AveragePrice           float64   `json:"averagePrice"`
	Count                  int       `json:"count"`
	FloorPriceMainCurrency float64   `json:"floorPriceMc"`
	FloorPriceUsd          float64   `json:"floorPriceUsd"`
	MainCurrencyID         string    `json:"mainCurrencyId"`
	MarketcapMainCurrency  float64   `json:"marketcapMc"`
	MarketcapUsd           float64   `json:"marketcapUsd"`
	OneDayAveragePrice     float64   `json:"oneDayAveragePrice"`
	OneDaySales            int       `json:"oneDaySales"`
	OwnersCount            int       `json:"ownersCount"`
	SevenDayAveragePrice   float64   `json:"sevenDayAveragePrice"`
	SevenDaySales          int       `json:"sevenDaySales"`
	ThirtyDayAveragePrice  float64   `json:"thirtyDayAveragePrice"`
	ThirtyDaySales         int       `json:"thirtyDaySales"`
	ThirtyDayVolume        float64   `json:"thirtyDayVolume"`
	TotalSales             int       `json:"totalSales"`
	TotalSupply            int       `json:"totalSupply"`
	TotalVolume            float64   `json:"totalVolume"`
	VolumeMainCurrency24h  float64   `json:"volumeMc24h"`
	VolumeMainCurrency7d   float64   `json:"volumeMc7d"`
	VolumeUsd24h           float64   `json:"volumeUsd24h"`
	FloorPriceChange24h    float64   `json:"floorPriceChange24h"`
	FloorPriceChange7d     float64   `json:"floorPriceChange7d"`
	MarketcapChange24h     float64   `json:"marketcapChange24h"`
	MarketcapChange7d      float64   `json:"marketcapChange7d"`
	VolumeChange24h        float64   `json:"volumeChange24h"`
	VolumeChange7d         float64   `json:"volumeChange7d"`
	OwnersCountChange24h   float64   `json:"ownersCountChange24h"`
	OwnersCountChange7d    float64   `json:"ownersCountChange7d"`
	SalesInProfit          int       `json:"salesInProfit"`
	SalesInProfitChange24h int       `json:"salesInProfitChange24h"`
	SalesInProfitChange7d  int       `json:"salesInProfitChange7d"`
	OneDayChange           float64   `json:"oneDayChange"`
	TransactionsUpdateDate time.Time `json:"transactionsUpdateDate"`
	SevenDayChange         float64   `json:"sevenDayChange"`
	ThirtyDayChange        float64   `json:"thirtyDayChange"`
	ListedCount            int       `json:"listedCount"`
	Show                   bool      `json:"show"`
	CreatorFee             float64   `json:"creatorFee"`
	Volume                 float64   `json:"volume"`
	CreatedDate            time.Time `json:"createdDate"`
}

type NftAsset struct {
	Address    string `json:"address"`
	Blockchain string `json:"blockchain"`
	TokenID    string `json:"tokenId"`
	Attributes []struct {
		Key           string    `json:"key"`
		Kind          string    `json:"kind"`
		Value         string    `json:"value"`
		TokenCount    int       `json:"tokenCount"`
		OnSaleCount   int       `json:"onSaleCount"`
		FloorAskPrice *float64  `json:"floorAskPrice"`
		TopBidValue   *float64  `json:"topBidValue"`
		CreatedAt     time.Time `json:"createdAt"`
	} `json:"attributes"`
	CollectionID  string    `json:"collectionId"`
	Name          string    `json:"name"`
	PreviewURL    string    `json:"previewUrl"`
	RarityRank    int       `json:"rarityRank"`
	RarityScore   float64   `json:"rarityScore"`
	Source        string    `json:"source"`
	Standard      string    `json:"standard"`
	URL           string    `json:"url"`
	LastSaleDate  time.Time `json:"lastSaleDate"`
	LastSalePrice float64   `json:"lastSalePrice"`
	ListSource    string    `json:"listSource"`
	ListPrice     float64   `json:"listPrice"`
}

type meta struct {
	Page            int  `json:"page"`
	Limit           int  `json:"limit"`
	ItemCount       int  `json:"itemCount"`
	PageCount       int  `json:"pageCount"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	HasNextPage     bool `json:"hasNextPage"`
}

type metaShort struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type coin struct {
	Rank           int     `json:"rank"`
	Identifier     string  `json:"identifier"`
	Symbol         string  `json:"symbol"`
	Name           string  `json:"name"`
	Icon           string  `json:"icon"`
	PriceChange24h float64 `json:"priceChange24h"`
	PriceChange1h  float64 `json:"priceChange1h"`
	PriceChange7d  float64 `json:"priceChange7d"`
	Volume         float64 `json:"volume"`
	IsFake         bool    `json:"isFake"`
	IsFiat         bool    `json:"isFiat"`
}
