package entity

import (
	"errors"
	"fmt"
	"time"
)

var (
	// репо
	ErrorDBConnect   = errors.New("cant connect to DB")
	ErrorInsertDB    = errors.New("insert failed")
	ErrorOrderExists = errors.New("order with such UID already exists")
	// для контроллера
	ErrInvalidInput    = errors.New("invalid input")
	ErrInternal        = errors.New("internal error")
	ErrAlreadyExists   = errors.New("already exists")
	ErrorOrderNotFound = errors.New("order not found")
	ErrorQueryFailed   = errors.New("request failed")
)

type OrderInfo struct {
	OrderUID          string       `json:"order_uid"`
	TrackNumber       string       `json:"track_number"`
	Entry             string       `json:"entry"`
	Delivery          DeliveryInfo `json:"delivery"`
	Payment           PaymentInfo  `json:"payment"`
	Items             []ItemInfo   `json:"items"`
	Locale            string       `json:"locale"`
	InternalSignature string       `json:"internal_signature"`
	CustomerID        string       `json:"customer_id"`
	DeliveryService   string       `json:"delivery_service"`
	ShardKey          string       `json:"shardkey"`
	SmID              int          `json:"sm_id"`
	DateCreated       time.Time    `json:"date_created"`
	OofShard          string       `json:"oof_shard"`
}

// validated

type DeliveryInfo struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type PaymentInfo struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int64  `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int64  `json:"delivery_cost"`
	GoodsTotal   int64  `json:"goods_total"`
	CustomFee    int64  `json:"custom_fee"`
}

type ItemInfo struct {
	ChrtID      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int64  `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int64  `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int64  `json:"status"`
}

type RequestIDKey struct{}

func (o *OrderInfo) ValidateOrder() error {
	if o == nil {
		return errors.New("nil order")
	}

	if o.OrderUID == "" {
		return errors.New("empty order_uid")
	}
	if o.TrackNumber == "" {
		return errors.New("empty track_number")
	}
	if o.Entry == "" {
		return errors.New("empty entry")
	}
	if o.Locale == "" {
		return errors.New("empty locale")
	}
	if o.CustomerID == "" {
		return errors.New("empty customer_id")
	}
	if o.DeliveryService == "" {
		return errors.New("empty delivery_service")
	}
	if o.DateCreated.IsZero() {
		return errors.New("empty date_created")
	}

	// Delivery
	if o.Delivery.Name == "" {
		return errors.New("empty delivery.name")
	}
	if o.Delivery.Phone == "" {
		return errors.New("empty delivery.phone")
	}
	if o.Delivery.City == "" {
		return errors.New("empty delivery.city")
	}
	if o.Delivery.Address == "" {
		return errors.New("empty delivery.address")
	}
	if o.Delivery.Email == "" {
		return errors.New("empty delivery.email")
	}

	// Payment
	if o.Payment.Transaction == "" {
		return errors.New("empty payment.transaction")
	}
	if o.Payment.Currency == "" {
		return errors.New("empty payment.currency")
	}
	if o.Payment.Provider == "" {
		return errors.New("empty payment.provider")
	}
	if o.Payment.Amount <= 0 {
		return errors.New("invalid payment.amount")
	}
	if o.Payment.PaymentDT <= 0 {
		return errors.New("invalid payment.payment_dt")
	}
	if o.Payment.Bank == "" {
		return errors.New("empty payment.bank")
	}

	// Items
	if len(o.Items) == 0 {
		return errors.New("empty items")
	}
	for i, it := range o.Items {
		if it.ChrtID == 0 {
			return fmt.Errorf("item[%d]: empty chrt_id", i)
		}
		if it.TrackNumber == "" {
			return fmt.Errorf("item[%d]: empty track_number", i)
		}
		if it.Name == "" {
			return fmt.Errorf("item[%d]: empty name", i)
		}
		if it.Price <= 0 {
			return fmt.Errorf("item[%d]: invalid price", i)
		}
		if it.TotalPrice <= 0 {
			return fmt.Errorf("item[%d]: invalid total_price", i)
		}
	}

	return nil
}
