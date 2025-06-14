package repositories

import (
	"fmt"
	"github.com/ewik2k21/grpc-hard/internal/models"
	"log/slog"
)

type ISpotInstrumentRepository interface {
	GetAllMarkets() (map[string]models.Market, error)
}

type SpotInstrumentRepository struct {
	markets map[string]models.Market
	logger  *slog.Logger
}

func NewSpotInstrumentRepository(logger *slog.Logger) *SpotInstrumentRepository {
	return &SpotInstrumentRepository{
		markets: make(map[string]models.Market),
		logger:  logger,
	}
}

func (r *SpotInstrumentRepository) GetAllMarkets() (map[string]models.Market, error) {
	if len(r.markets) == 0 {
		r.logger.Error("no markets in memory")
		return nil, fmt.Errorf("zero markets")
	}
	return r.markets, nil
}
