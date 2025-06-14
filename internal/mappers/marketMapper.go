package mappers

import (
	"github.com/ewik2k21/grpc-hard/internal/models"
	"github.com/google/uuid"
)
import spot "github.com/ewik2k21/grpc-hard/pkg/spot_instrument_service_v1"

func MapMarketToProto(m models.Market) *spot.Market {
	return &spot.Market{
		Id:   m.ID.String(),
		Name: m.Name,
	}
}

func MapProtoToMarkets(resp *spot.ViewMarketsResponse) ([]*models.Market, error) {
	markets := resp.GetMarkets()
	res := make([]*models.Market, 0, len(markets))
	for _, market := range markets {
		marketId, err := uuid.Parse(market.Id)
		if err != nil {
			return nil, err
		}
		res = append(res, &models.Market{
			ID:        marketId,
			Name:      market.Name,
			Enabled:   true,
			DeletedAt: nil,
		})
	}
	return res, nil
}
