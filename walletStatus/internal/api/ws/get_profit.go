package wsserver

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
)

type getProfit struct {
	Ticker        string  `json:"ticker"`
	Profit        float64 `json:"profit"`
	ProfitPercent float64 `json:"profit_percent"`
	Balance       float64 `json:"balance"`
}

type getProfitByTickersRes struct {
	Profits []getProfit `json:"profits"`
	Revenue float64     `json:"revenue"`
}

func (srv *Server) handleProfitWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ch, errCh := srv.stremService.StreamProfitByBit()

	for {
		select {
		case profit := <-ch:
			profitRes := make([]getProfit, 0, len(profit))

			var revenue float64

			for ticker, data := range profit {
				revenue += data.Profit

				profitRes = append(profitRes, getProfit{
					Ticker:        ticker,
					Profit:        data.Profit,
					ProfitPercent: data.ProfitPercent,
					Balance:       data.Balance,
				})
			}

			sort.Slice(profitRes, func(i, j int) bool {
				return profitRes[i].Balance > profitRes[j].Balance
			})

			res := getProfitByTickersRes{
				Profits: profitRes,
				Revenue: revenue,
			}

			payload, _ := json.Marshal(res)
			conn.WriteMessage(websocket.TextMessage, payload)
		case err := <-errCh:
			conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"`+err.Error()+`"}`))

		}
	}
}
