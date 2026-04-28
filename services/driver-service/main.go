package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9092"

func main() {
	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service := NewService()
	// starting the gRPC server
	grpcServer := grpcserver.NewServer()
	NewGrpcHandler(grpcServer, service)
	log.Printf("Starting gRPC server Driver service on %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done

	log.Println("Shutting down gRPC server...")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-stopped:
		log.Println("gRPC server stopped gracefully")
	case <-ctx.Done():
		log.Println("Timeout exceeded, forcing server stop")
		grpcServer.Stop()
	}
}
