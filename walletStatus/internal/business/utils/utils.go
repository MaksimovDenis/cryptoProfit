package utils

import (
	"strings"
	"walletStatus/internal/domain"
)

func ConvTickersToStrBinance(priceByTicker map[string]domain.Stock) string {
	var sb strings.Builder

	sb.WriteString("[")

	idx := 0

	for ticker := range priceByTicker {
		if idx > 0 {
			sb.WriteString(",")
		}

		sb.WriteString("\"" + ticker + "\"")

		idx++
	}

	sb.WriteString("]")
	finalString := sb.String()

	return finalString
}

func ConvTickersToStrByBit(priceByTicker map[string]domain.Stock) []string {
	const chunkSize = 10
	result := []string{}

	chunk := []string{}
	for ticker := range priceByTicker {
		chunk = append(chunk, "\"tickers."+ticker+"\"")

		if len(chunk) == chunkSize {
			result = append(result, "["+strings.Join(chunk, ",")+"]")
			chunk = []string{}
		}
	}

	if len(chunk) > 0 {
		result = append(result, "["+strings.Join(chunk, ",")+"]")
	}

	return result
}
