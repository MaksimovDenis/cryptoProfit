package api

import (
	"encoding/json"
	"net/http"
	"sort"
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

func (s *Server) GetProfitByTickers(w http.ResponseWriter, r *http.Request) {
	profit, err := s.walletStatus.GetProfitByTickers(r.Context())
	if err != nil {
		makeErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

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

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		makeErrorResponse(w, err, http.StatusInternalServerError)

		return
	}
}
