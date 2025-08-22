package usecase

import (
	"context"
	"errors"

	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"go.uber.org/zap"
)

func (u *UsecaseLayer) GetOrderInfo(ctx context.Context, orderUID string) (*entity.OrderResponse, error) {
	// 1) забираем request_id
	reqID, _ := ctx.Value(entity.RequestIDKey{}).(string)

	// 2) оборачиваем логгер
	logger := u.log.With(zap.String("func", "GetOrderInfo"))
	if reqID != "" {
		logger = logger.With(zap.String("request_id", reqID))
	}

	if orderUID == "" {
		logger.Warn("empty order_uid")
		return nil, entity.ErrInvalidInput
	}

	order, err := u.db.GetOrderByUID(ctx, orderUID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrorOrderNotFound):
			logger.Info("order not found")
			return nil, entity.ErrorOrderNotFound // контроллер вернёт 404
		case errors.Is(err, entity.ErrorQueryFailed):
			// если контекст уже отменён/истёк — вернём 504 на уровне контроллера
			if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.Canceled) {
				logger.Error("query failed due to ctx deadline/cancel", zap.Error(ctx.Err()))
				return nil, entity.ErrInternal
			}
			logger.Error("query failed", zap.Error(err))
			return nil, entity.ErrInternal // 500
		default:
			if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.Canceled) {
				logger.Error("unexpected repo error with ctx cancel/deadline", zap.Error(err), zap.Error(ctx.Err()))
				return nil, entity.ErrInternal
			}
			logger.Error("unexpected repo error", zap.Error(err))
			return nil, entity.ErrInternal
		}
	}

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
