package entity

import (
	"time"
)

type OrderResponse struct {
	OrderUID    string         `json:"order_uid"`
	DateCreated time.Time      `json:"date_created"`
	Locale      string         `json:"locale"`
	Logistics   LogisticsInfo  `json:"logistics"`
	Delivery    DeliveryPublic `json:"delivery"`
	Payment     PaymentPublic  `json:"payment"`
	Items       []ItemPublic   `json:"items"`
}

type LogisticsInfo struct {
	TrackNumber     string `json:"track_number"`
	DeliveryService string `json:"delivery_service"`
}

type DeliveryPublic struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

type PaymentPublic struct {
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	DeliveryCost int64  `json:"delivery_cost"`
	GoodsTotal   int64  `json:"goods_total"`
}

type ItemPublic struct {
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	Size       string `json:"size"`
	Price      int64  `json:"price"`
	TotalPrice int64  `json:"total_price"`
	Status     int64  `json:"status"`
}
