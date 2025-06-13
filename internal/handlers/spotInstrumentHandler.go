package handlers

import (
	spotInstrument "github.com/ewik2k21/grpc-hard/pkg/spot_instrument_service_v1"
)

type SpotInstrumentHandler struct {
	spotInstrument.UnimplementedSpotInstrumentServiceServer
	//todo add service
}

func NewSpotInstrumentHandler() *SpotInstrumentHandler {
	return &SpotInstrumentHandler{}
}
