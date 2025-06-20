package handlers

import (
	"context"
	"github.com/ewik2k21/grpc-hard/internal/services"
	order "github.com/ewik2k21/grpc-hard/pkg/order_service_v1"
	spotInstrument "github.com/ewik2k21/grpc-hard/pkg/spot_instrument_service_v1"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type OrderHandler struct {
	order.UnimplementedOrderServiceServer
	Client  spotInstrument.SpotInstrumentServiceClient
	service services.OrderService
	logger  *slog.Logger
	//todo add service
}

func NewOrderHandler(
	logger *slog.Logger,
	service *services.OrderService,
	redisClient *redis.Client,
) *OrderHandler {
	return &OrderHandler{
		logger:  logger,
		service: *service,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {

	userRole := request.GetUserRole()
	resp, err := h.Client.ViewMarkets(
		context.Background(),
		&spotInstrument.ViewMarketsRequest{
			UserRole: userRole,
		})
	if err != nil {
		return nil, err
	}

	orderId, status, err := h.service.CreateOrder(resp, request)
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		OrderId: orderId,
		Status:  *status,
	}, nil

}

func (h *OrderHandler) GetOrderStatus(ctx context.Context, req *order.GetOrderStatusRequest) (*order.GetOrderStatusResponse, error) {

	userId := req.GetUserId()
	orderId := req.GetOrderId()

	status, err := h.service.GetOrderStatus(userId, orderId)
	if err != nil {
		return nil, err
	}

	return &order.GetOrderStatusResponse{
		Status: *status,
	}, nil
}
