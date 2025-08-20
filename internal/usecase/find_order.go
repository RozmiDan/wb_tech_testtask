package usecase

import (
	"context"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"go.uber.org/zap"
)

type RepoLayer interface {
	GetOrderByUID(ctx context.Context, orderUID string) (*entity.OrderInfo, error)
}

type UsecaseLayer struct {
	log *zap.Logger
	db  RepoLayer
}

func New(logger *zap.Logger, dbLayer RepoLayer) *UsecaseLayer {
	return &UsecaseLayer{
		log: logger.With(zap.String("layer", "Usecase")),
		db:  dbLayer,
	}
}

func (u *UsecaseLayer) GetOrderInfo(ctx context.Context, orderUID string) (*entity.OrderResponse, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := u.log.With(zap.String("func", "GetOrderInfo"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	order, err := u.db.GetOrderByUID(ctx, orderUID)
	if err != nil {
		// TODO
	}

	logger.Info("successfuly found order")

	resOrd := mapOrderToResponse(order)

	logger.Info("succsessfuly found order")
	return resOrd, nil
}

func mapOrderToResponse(order *entity.OrderInfo) *entity.OrderResponse {
	items := make([]entity.ItemPublic, 0, len(order.Items))
	for _, it := range order.Items {
		items = append(items, entity.ItemPublic{
			Name:       it.Name,
			Brand:      it.Brand,
			Size:       it.Size,
			Price:      it.Price,
			TotalPrice: it.TotalPrice,
			Status:     it.Status,
		})
	}

	return &entity.OrderResponse{
		OrderUID:    order.OrderUID,
		DateCreated: order.DateCreated,
		Locale:      order.Locale,
		Logistics: entity.LogisticsInfo{
			TrackNumber:     order.TrackNumber,
			DeliveryService: order.DeliveryService,
		},
		Delivery: entity.DeliveryPublic{
			Name:    order.Delivery.Name,
			City:    order.Delivery.City,
			Region:  order.Delivery.Region,
			Address: order.Delivery.Address,
			Email:   order.Delivery.Email,
			Phone:   order.Delivery.Phone,
		},
		Payment: entity.PaymentPublic{
			Amount:       order.Payment.Amount,
			Currency:     order.Payment.Currency,
			DeliveryCost: order.Payment.DeliveryCost,
			GoodsTotal:   order.Payment.GoodsTotal,
		},
		Items: items,
	}
}
