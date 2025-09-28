package wsserver

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
)

type getPriceByTickersRes struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:5173"
	},
}

func (srv *Server) handlePriceWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ch, errCh := srv.stremService.StreamPricesByBit()

	for {
		select {
		case data := <-ch:
			res := make([]getPriceByTickersRes, 0, len(data))

			for symbol, price := range data {
				res = append(res, getPriceByTickersRes{
					Symbol: symbol,
					Price:  price,
				})
			}

			sort.Slice(res, func(i, j int) bool {
				return res[i].Symbol > res[j].Symbol
			})

			payload, _ := json.Marshal(res)
			conn.WriteMessage(websocket.TextMessage, payload)
		case err := <-errCh:
			conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"`+err.Error()+`"}`))

		}

	}
}
