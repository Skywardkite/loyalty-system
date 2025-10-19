package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type OrderInfo struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float32 `json:"accrual"`
}

func (c *Client) GetOrderAccrual(ctx context.Context, orderNum string) (*OrderInfo, int, int, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNum)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := 60 // дефолтное значение
		if hdr := resp.Header.Get("Retry-After"); hdr != "" {
			if sec, err := strconv.Atoi(strings.TrimSpace(hdr)); err == nil {
				retryAfter = sec
			}
		}
		return nil, resp.StatusCode, retryAfter, nil
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, resp.StatusCode, 0, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, 0, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var info OrderInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, resp.StatusCode, 0, err
	}

	return &info, resp.StatusCode, 0, nil
}
