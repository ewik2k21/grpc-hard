package handlers

import (
	order "github.com/ewik2k21/grpc-hard/pkg/order_service_v1"
	spotInstrument "github.com/ewik2k21/grpc-hard/pkg/spot_instrument_service_v1"
)

type OrderHandler struct {
	order.UnimplementedOrderServiceServer
	Client spotInstrument.SpotInstrumentServiceClient
	//todo add service
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}
