package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9092"

func main() {
	rabbitmqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service := NewService()
	rabbitmq, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()
	// starting the gRPC server
	grpcServer := grpcserver.NewServer()
	NewGrpcHandler(grpcServer, service)
	log.Printf("Starting gRPC server Driver service on %s", lis.Addr().String())

	consumer := NewTripConsumer(rabbitmq)
	go func() {
		if err := consumer.Listen(); err != nil {
			log.Fatalf("Failed to listen to the message: %v", err)
		}
	}()
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
