package api

import (
	"context"
	"net/http"
	"walletStatus/internal/domain"
	middleware "walletStatus/internal/middlware"
)

type walletStatusService interface {
	GetPriceByTickers(ctx context.Context) (map[string]float64, error)
	GetProfitByTickers(ctx context.Context) (map[string]domain.Profit, error)
}

type Server struct {
	walletStatus walletStatusService
}

func New(walletStatus walletStatusService) *Server {
	return &Server{
		walletStatus: walletStatus,
	}
}

func (s *Server) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/profit", s.GetProfitByTickers)
	mux.HandleFunc("GET /api/prices", s.GetPriceByTickers)

	h := middleware.WithCORS(mux)

	return h
}
