package server

import (
	"context"
	"github.com/ewik2k21/grpc-hard/config"
	"github.com/ewik2k21/grpc-hard/internal/handlers"
	"github.com/ewik2k21/grpc-hard/internal/interceptors/loggerInterceptor"
	x_request_id "github.com/ewik2k21/grpc-hard/internal/interceptors/x-request-id"
	"github.com/ewik2k21/grpc-hard/internal/repositories"
	"github.com/ewik2k21/grpc-hard/internal/services"
	order "github.com/ewik2k21/grpc-hard/pkg/order_service_v1"
	spotInstrument "github.com/ewik2k21/grpc-hard/pkg/spot_instrument_service_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const bufSize = 1024 * 1024

func Execute(logger *slog.Logger) {
	config.LoadConfig()
	lis := bufconn.Listen(bufSize)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			x_request_id.RequestIDInterceptor(),
			loggerInterceptor.LoggerRequestInterceptor(logger),
		),
	)

	orderRepo := repositories.NewOrderRepository(logger)
	orderService := services.NewOrderService(*orderRepo, logger)
	orderHandler := handlers.NewOrderHandler(logger, *orderService)

	spotInstrumentRepo := repositories.NewSpotInstrumentRepository(logger)
	spotInstrumentService := services.NewSpotInstrumentService(*spotInstrumentRepo, logger)
	spotInstrumentHandler := handlers.NewSpotInstrumentHandler(*spotInstrumentService, logger)
	order.RegisterOrderServiceServer(grpcServer, orderHandler)
	spotInstrument.RegisterSpotInstrumentServiceServer(grpcServer, spotInstrumentHandler)

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("failed to create conn: %v", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close()

	orderHandler.Client = spotInstrument.NewSpotInstrumentServiceClient(conn)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("start grpc server")

		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("grpc server not started", slog.String("error", err.Error()))
		}
	}()

	lis2, err := net.Listen("tcp", os.Getenv(config.GrpcPort))
	if err != nil {
		logger.Error("failed to lis tcp", slog.String("error", err.Error()))
		os.Exit(1)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("start tcp server")
		if err := grpcServer.Serve(lis2); err != nil {
			logger.Error("tcp server not started", slog.String("error", err.Error()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("received shutdown signal, start graceful shutdown")

	grpcServer.GracefulStop()
	logger.Info("server stopped")

	wg.Wait()
	logger.Info("all stopped")

}
