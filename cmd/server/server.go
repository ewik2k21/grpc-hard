package server

import (
	"context"
	"github.com/ewik2k21/grpc-hard/config"
	"github.com/ewik2k21/grpc-hard/internal/handlers"
	"github.com/ewik2k21/grpc-hard/internal/interceptors"
	"github.com/ewik2k21/grpc-hard/internal/repositories"
	"github.com/ewik2k21/grpc-hard/internal/services"
	order "github.com/ewik2k21/grpc-hard/pkg/order_service_v1"
	spotInstrument "github.com/ewik2k21/grpc-hard/pkg/spot_instrument_service_v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const bufSize = 1024 * 1024

func Execute(logger *slog.Logger) {
	ctx := context.Background()
	err := config.LoadConfig()
	if err != nil {
		logger.Error("failed load env file", slog.String("error", err.Error()))
		os.Exit(1)
	}
	redisClient, err := config.InitRedis(ctx)
	if err != nil {
		logger.Error("failed init redis", slog.String("error", err.Error()))
		os.Exit(1)
	}
	lis := bufconn.Listen(bufSize)
	g, errCtx := errgroup.WithContext(ctx)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestIDInterceptor(),
			interceptors.LoggerRequestInterceptor(logger),
			interceptors.UnaryPanicRecoveryInterceptor(logger),
			interceptors.PrometheusInterceptor(),
		),
	)

	orderRepo := repositories.NewOrderRepository(logger)
	orderService := services.NewOrderService(*orderRepo, logger)
	orderHandler := handlers.NewOrderHandler(logger, orderService, redisClient)

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
	g.Go(func() error {
		defer wg.Done()
		logger.Info("start grpc server")

		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("grpc server not started", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	lis2, err := net.Listen("tcp", os.Getenv(config.GrpcPort))
	if err != nil {
		logger.Error("failed to lis tcp", slog.String("error", err.Error()))
		os.Exit(1)
	}

	wg.Add(1)
	g.Go(func() error {
		defer wg.Done()
		logger.Info("start tcp server")
		if err := grpcServer.Serve(lis2); err != nil {
			logger.Error("tcp server not started", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	//server for metrics
	metricsServer := &http.Server{
		Addr: ":2112",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				promhttp.Handler().ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		}),
	}

	wg.Add(1)
	g.Go(func() error {
		defer wg.Done()
		logger.Info("metrics endpoint start on :2112")
		http.Handle("/metrics", promhttp.Handler())
		if err := metricsServer.ListenAndServe(); err != nil {
			logger.Error("metrics endpoint failed", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("received shutdown signal, start graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(errCtx, 5*time.Second)
	defer cancel()
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("metrics server shutdown failed", slog.String("error", err.Error()))
	}

	grpcServer.GracefulStop()
	logger.Info("server stopped")

	wg.Wait()
	logger.Info("all stopped")

}
