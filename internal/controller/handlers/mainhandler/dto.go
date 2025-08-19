package mainhandler

import "github.com/RozmiDan/wb_tech_testtask/internal/entity"

type GetOrderResponse struct {
	Order entity.OrderResponse
}

// APIError — единая структура описания ошибки
type APIError struct {
	Code    string `json:"code"`    // машинно-читаемый код ошибки
	Message string `json:"message"` // человеко-читаемое сообщение
}
