package services

import (
	"fmt"
	"github.com/ewik2k21/grpc-hard/internal/mappers"
	"github.com/ewik2k21/grpc-hard/internal/repositories"
	order "github.com/ewik2k21/grpc-hard/pkg/order_service_v1"
	pkg "github.com/ewik2k21/grpc-hard/pkg/spot_instrument_service_v1"
	"github.com/google/uuid"
	"log/slog"
)

type OrderService struct {
	repo   repositories.OrderRepository
	logger *slog.Logger
}

func NewOrderService(repo repositories.OrderRepository, logger *slog.Logger) *OrderService {
	return &OrderService{
		repo:   repo,
		logger: logger,
	}
}

func (s *OrderService) CreateOrder(resp *pkg.ViewMarketsResponse, request *order.CreateOrderRequest) (string, *order.Status, error) {

	markets, err := mappers.MapProtoToMarkets(resp)
	if err != nil {
		s.logger.Error("failed mapping proto to markets", slog.String("error", err.Error()))
		return "", nil, err
	}

	marketId := request.GetMarketId()
	var ok bool = false

	mapOrder, err := mappers.MapProtoToOrder(request)
	if err != nil {
		s.logger.Error("failed mapping proto to order", slog.String("error", err.Error()))
		return "", nil, err
	}

	var orderId *uuid.UUID
	var status *order.Status
	for _, market := range markets {
		if market.ID.String() == marketId {
			ok = true
			orderId, status, err = s.repo.CreateOrder(mapOrder)
		}
	}

	if !ok {
		return "", nil, fmt.Errorf("needed market not found")
	}

	return orderId.String(), status, nil

}
