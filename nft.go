package coinstats

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type NftTrendingFilter struct {
	Page  int
	Limit int
}

type NftTrendingResult struct {
	Meta meta            `json:"meta"`
	Data []NftCollection `json:"data"`
}

// NftTrending return the most popular NFT collections right now.
func (c *Client) NftTrending(ctx context.Context, filter NftTrendingFilter) (NftTrendingResult, error) {
	q := url.Values{}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/nft/trending?"+q.Encode(), nil)
	res, err := do[NftTrendingResult](c, req)
	return res, err
}

type NftsByWalletFilter struct {
	Page  int
	Limit int
}

type NftsByWalletResult struct {
	Meta meta `json:"meta"`
	Data []struct {
		Name               string  `json:"name"`
		Logo               string  `json:"logo"`
		Address            string  `json:"address"`
		TotalFloorPrice    float64 `json:"totalFloorPrice"`
		TotalLastSalePrice float64 `json:"totalLastSalePrice"`
		ID                 string  `json:"id"`
		AssetsCount        int     `json:"assetsCount"`
		Assets             []struct {
			PreviewImg string `json:"previewImg"`
		} `json:"assets"`
		FloorPrice float64 `json:"floorPrice"`
	} `json:"data"`
}

// NftsByWallet return list of NFT assets owned by a wallet address.
func (c *Client) NftsByWallet(ctx context.Context, addr string, filter NftsByWalletFilter) (NftsByWalletResult, error) {
	q := url.Values{}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		host+"/nft/wallet/"+url.QueryEscape(addr)+"/assets?"+q.Encode(), nil)
	res, err := do[NftsByWalletResult](c, req)
	return res, err
}

// NftCollectionByAddress return detailed information about an NFT collection using collectionAddress.
func (c *Client) NftCollectionByAddress(ctx context.Context, addr string) (NftCollection, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/nft/collection/"+url.QueryEscape(addr), nil)
	res, err := do[NftCollection](c, req)
	return res, err
}

type NftCollectionAssetFilter struct {
	Page   int
	Limit  int
	Listed bool
}

type NftCollectionAssetsByAddressResult struct {
	Meta meta       `json:"meta"`
	Data []NftAsset `json:"data"`
}

// NftCollectionAssetsByAddress return the list of NFT assets associated with NFT Collection by collectionAddress.
func (c *Client) NftCollectionAssetsByAddress(ctx context.Context, addr string, filter NftCollectionAssetFilter) (NftCollectionAssetsByAddressResult, error) {
	q := url.Values{}
	if filter.Page > 0 {
		q.Add("page", strconv.Itoa(filter.Page))
	}
	if filter.Limit > 0 {
		q.Add("limit", strconv.Itoa(filter.Limit))
	}
	t := "all"
	if filter.Listed {
		t = "listed"
	}
	q.Add("type", t)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		host+"/nft/"+url.QueryEscape(addr)+"/assets?"+q.Encode(), nil)
	res, err := do[NftCollectionAssetsByAddressResult](c, req)
	return res, err
}

// NftCollectionAssetByToken return detailed information about a specific NFT asset.
func (c *Client) NftCollectionAssetByToken(ctx context.Context, addr, tokenID string) (NftAsset, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		host+"/nft/"+url.QueryEscape(addr)+"/asset/"+url.QueryEscape(tokenID), nil)
	res, err := do[NftAsset](c, req)
	return res, err
}
