package domain

import "errors"

var (
	ErrTickersUnique   = errors.New("Каждый тикер должен быть уникален")
	ErrTickersNotFound = errors.New("Не заданы тикеры")
	ErrPricesCount     = errors.New("Количество тикеров не соответсвует количеству средних ценовых позиций")

	ErrFailedToParseAddr = errors.New("failed to parse address")
)
