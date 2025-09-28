package wsserver

import (
	"net/http"
	"walletStatus/internal/domain"
)

type stremService interface {
	StreamPricesByBit() (<-chan map[string]float64, <-chan error)
	StreamProfitByBit() (chan map[string]domain.Profit, <-chan error)
}

type Server struct {
	stremService stremService
}

func New(stremService stremService) *Server {
	return &Server{
		stremService: stremService,
	}
}

func (s *Server) Run(addr string) error {
	http.HandleFunc("/ws/price", s.handlePriceWS)
	http.HandleFunc("/ws/profit", s.handleProfitWS)
	return http.ListenAndServe(addr, nil)
}
