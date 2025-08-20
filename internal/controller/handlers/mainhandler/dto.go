package mainhandler

import "github.com/RozmiDan/wb_tech_testtask/internal/entity"

type GetOrderResponse struct {
	Order entity.OrderResponse
}

type APIError struct {
	Message string `json:"message"` 
}
