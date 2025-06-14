package handlers

import (
	"context"
	"github.com/ewik2k21/grpc-hard/internal/services"
	spotInstrument "github.com/ewik2k21/grpc-hard/pkg/spot_instrument_service_v1"
	"log/slog"
)

type SpotInstrumentHandler struct {
	spotInstrument.UnimplementedSpotInstrumentServiceServer
	service services.SpotInstrumentService
	logger  *slog.Logger
}

func NewSpotInstrumentHandler(service services.SpotInstrumentService, logger *slog.Logger) *SpotInstrumentHandler {
	return &SpotInstrumentHandler{service: service, logger: logger}
}

func (h *SpotInstrumentHandler) ViewMarkets(ctx context.Context, req *spotInstrument.ViewMarketsRequest) (res *spotInstrument.ViewMarketsResponse, err error) {
	userRole := req.GetUserRole()

	markets, err := h.service.GetAllMarkets(userRole)
	if err != nil {
		return nil, err
	}

	return &spotInstrument.ViewMarketsResponse{
		Markets: markets,
	}, nil
}
