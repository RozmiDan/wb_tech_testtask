package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RozmiDan/wb_tech_testtask/internal/config"
	"github.com/RozmiDan/wb_tech_testtask/internal/entity"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type OrderHandler interface {
	AddOrderInfo(ctx context.Context, order *entity.OrderInfo) error
}

type Consumer struct {
	reader  *kafka.Reader
	handler OrderHandler
	logger  *zap.Logger
}

func NewConsumer(cfg *config.Config, handler OrderHandler, logger *zap.Logger) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.KafkaBrokers,
		GroupID:  cfg.KafkaGroupID,
		Topic:    cfg.KafkaTopic,
		MinBytes: cfg.KafkaMinBytes,
		MaxBytes: cfg.KafkaMaxBytes,
	})
	return &Consumer{
		reader:  r,
		handler: handler,
		logger:  logger.With(zap.String("component", "kafka_consumer"), zap.String("topic", cfg.KafkaTopic)),
	}
}

func (c *Consumer) Start(ctx context.Context, cfg *config.Config) {
	c.logger.Info("starting kafka consumer loop")
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				c.logger.Info("context done, exiting consumer")
				return
			}
			c.logger.Error("fetch message failed", zap.Error(err))
			continue
		}

		// формируем request_id и контекст на обработку одного сообщения
		reqID := uuid.NewString()

		ctxMsg := context.WithValue(ctx, entity.RequestIDKey{}, reqID)
		ctxMsg, cancel := context.WithTimeout(ctxMsg, cfg.KafkaMsgTimeout*time.Second)

		order := &entity.OrderInfo{}
		if err := json.Unmarshal(msg.Value, order); err != nil {
			c.logger.Warn("invalid message json, skipping",
				zap.String("request_id", reqID),
				zap.Int("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
				zap.Error(err),
			)
			_ = c.reader.CommitMessages(ctx, msg)
			cancel()
			continue
		}

		err = order.ValidateOrder()
		if err != nil {
			c.logger.Warn("invalid message payload, skipping",
				zap.String("request_id", reqID),
				zap.String("order_uid", order.OrderUID),
				zap.Error(err),
			)
			_ = c.reader.CommitMessages(ctx, msg)
			cancel()
			continue
		}

		if err := c.handler.AddOrderInfo(ctxMsg, order); err != nil {
			switch {
			case errors.Is(err, entity.ErrAlreadyExists):
				c.logger.Info("order already exists, committing",
					zap.String("request_id", reqID),
					zap.String("order_uid", order.OrderUID),
				)
				_ = c.reader.CommitMessages(ctx, msg)
			default:
				c.logger.Error("handler failed, will retry (no commit)",
					zap.String("request_id", reqID),
					zap.String("order_uid", order.OrderUID),
					zap.Error(err),
				)
			}
			cancel()
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Warn("commit failed",
				zap.String("request_id", reqID),
				zap.String("order_uid", order.OrderUID),
				zap.Error(err),
			)
		} else {
			c.logger.Info("order ingested",
				zap.String("request_id", reqID),
				zap.String("order_uid", order.OrderUID),
			)
		}
		cancel()
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
