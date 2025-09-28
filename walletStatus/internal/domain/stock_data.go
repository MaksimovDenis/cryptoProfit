package domain

type Stock struct {
	Ticker       string
	AveragePrice float64
	Balance      float64
}

type Profit struct {
	Profit        float64
	ProfitPercent float64
	Balance       float64
}
