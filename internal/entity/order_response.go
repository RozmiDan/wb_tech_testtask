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
	Name    string `json:"name_masked"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Address string `json:"address_masked"`
	Email   string `json:"email_masked"`
	Phone   string `json:"phone_masked"`
}

type PaymentPublic struct {
	Amount       int    `json:"amount"`
	Currency     string `json:"currency"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
}

type ItemPublic struct {
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	Size       string `json:"size"`
	Price      int    `json:"price"`
	TotalPrice int    `json:"total_price"`
	Status     int    `json:"status"`
}
