package bybitws

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"walletStatus/internal/infra/transport/ws"
)

type Client struct {
	ws    *ws.Client
	out   chan map[string]float64
	errCh chan error
}

type tickerMessageRes struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Data  struct {
		Symbol    string `json:"symbol"`
		LastPrice string `json:"lastPrice"`
		Bid1Price string `json:"bid1Price"`
		Ask1Price string `json:"ask1Price"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

func New(ctx context.Context, addr string) (*Client, error) {
	rawOut := make(chan []byte, 100)
	rawIn := make(chan []byte, 100)

	ws, err := ws.NewClient(ws.Opts{
		Ctx:  ctx,
		Addr: addr,
		In:   rawIn,
		Out:  rawOut,
	})
	if err != nil {
		return nil, fmt.Errorf("bybitws.New: failed to create bybit client: %w", err)
	}

	cl := &Client{
		ws:    ws,
		out:   make(chan map[string]float64, 100),
		errCh: make(chan error, 10),
	}

	go cl.listen(rawOut)
	ws.Start()

	return cl, nil
}

func (c *Client) SubscribeTickers(tickers string) (<-chan map[string]float64, chan error) {
	subMsg := fmt.Sprintf(`{
		"op": "subscribe",
		"args": %s
	}`, tickers)

	go func() { c.ws.In <- []byte(subMsg) }()

	return c.out, c.errCh
}

func (c *Client) Close() error {
	c.ws.Stop()
	return nil
}
func (c *Client) listen(raw <-chan []byte) {
	var payload tickerMessageRes
	res := make(map[string]float64)

	for msg := range raw {
		if err := json.Unmarshal(msg, &payload); err != nil {
			c.errCh <- err
			continue
		}

		if payload.Data.Symbol != "" {
			price, err := strconv.ParseFloat(payload.Data.LastPrice, 64)
			if err != nil {
				c.errCh <- fmt.Errorf("strconv.ParseFloat: %w", err)
				continue
			}

			res[payload.Data.Symbol] = price

			snapshot := make(map[string]float64, len(res))
			for k, v := range res {
				snapshot[k] = v
			}

			c.out <- snapshot
		}
	}
}
