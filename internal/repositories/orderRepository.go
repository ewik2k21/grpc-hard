package repositories

import (
	"github.com/ewik2k21/grpc-hard/internal/models"
	order "github.com/ewik2k21/grpc-hard/pkg/order_service_v1"
	"github.com/google/uuid"
	"log/slog"
)

type IOrderRepository interface {
	CreateOrder(order *models.Order) (uuid.UUID, *order.Status, error)
}

type OrderRepository struct {
	orders map[string]models.Order
	logger *slog.Logger
}

func NewOrderRepository(logger *slog.Logger) *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]models.Order),
		logger: logger,
	}
}

func (r *OrderRepository) CreateOrder(order *models.Order) (*uuid.UUID, *order.Status, error) {

	return nil, nil, nil
}
