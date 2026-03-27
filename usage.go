package coinstats

import (
	"context"
	"net/http"
)

type Credits struct {
	Total        int    `json:"totalCredits"`
	Used         int    `json:"usedCredits"`
	Remaining    int    `json:"remainingCredits"`
	Subscription string `json:"subscription"`
}

func (c *Client) CreditUsage(ctx context.Context) (Credits, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host+"/usage/credits", nil)
	credits, err := do[Credits](c, req)
	return credits, err
}
