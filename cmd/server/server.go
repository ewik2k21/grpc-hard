package server

import (
	"context"
	"github.com/ewik2k21/grpc-hard/config"
	"github.com/ewik2k21/grpc-hard/internal/handlers"
	order "github.com/ewik2k21/grpc-hard/pkg/order_service_v1"
	spotInstrument "github.com/ewik2k21/grpc-hard/pkg/spot_instrument_service_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const bufSize = 1024 * 1024

func Execute() {
	config.LoadConfig()
	lis := bufconn.Listen(bufSize)

	grpcServer := grpc.NewServer()
	orderHandler := handlers.NewOrderHandler()
	spotInstrumentHandler := handlers.NewSpotInstrumentHandler()
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
		log.Fatalf("failed to create conn: %v", err)
	}
	defer conn.Close()

	orderHandler.Client = spotInstrument.NewSpotInstrumentServiceClient(conn)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("start grpc server")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc server not started: %v", err)
		}
	}()

	lis2, err := net.Listen("tcp", os.Getenv(config.GrpcPort))
	if err != nil {
		log.Fatalf("faile to lis tcp: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("start tcp server")
		if err := grpcServer.Serve(lis2); err != nil {
			log.Fatalf("tcp server not started: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("received shutdown signal, start graceful shutdown")

	grpcServer.GracefulStop()
	log.Println("server stopped")

	wg.Wait()
	log.Println("all stopped")

}
