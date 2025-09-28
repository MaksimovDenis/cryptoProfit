package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"walletStatus/internal/domain"
)

type Client struct {
	httpClient http.Client
	address    string
}

type getPriceResponseDTO struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func New(httpClient http.Client, address string) *Client {
	return &Client{
		httpClient: httpClient,
		address:    address,
	}
}

func (cl *Client) GetPriceByTickers(ctx context.Context, tickers string) (map[string]float64, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		cl.address,
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	qParams := req.URL.Query()
	qParams.Add("symbols", tickers)
	req.URL.RawQuery = qParams.Encode()

	response, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient.Do %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance.GetPriceByTicker: %w", err)
	}

	var resDTO []getPriceResponseDTO
	if err := json.NewDecoder(response.Body).Decode(&resDTO); err != nil {
		return nil, fmt.Errorf("json.NewDecoder: %w", err)
	}

	res := make(map[string]float64, len(resDTO))

	for _, ticker := range resDTO {
		if _, ok := res[ticker.Symbol]; ok {
			return nil, domain.ErrTickersUnique
		}

		price, err := strconv.ParseFloat(ticker.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("strconv.ParseFloat: %w", err)
		}

		res[ticker.Symbol] = price
	}

	return res, nil
}
