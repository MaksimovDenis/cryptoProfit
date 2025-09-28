package api

import (
	"encoding/json"
	"net/http"
	"sort"
)

type getPriceByTickersRes struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

func (s *Server) GetPriceByTickers(w http.ResponseWriter, r *http.Request) {
	prices, err := s.walletStatus.GetPriceByTickers(r.Context())
	if err != nil {
		makeErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	res := make([]getPriceByTickersRes, 0, len(prices))

	for symbol, price := range prices {
		res = append(res, getPriceByTickersRes{
			Symbol: symbol,
			Price:  price,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Symbol > res[j].Symbol
	})

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		makeErrorResponse(w, err, http.StatusInternalServerError)

		return
	}
}
