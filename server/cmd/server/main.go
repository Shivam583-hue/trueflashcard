package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shivam583-hue/trueflashcard/server/internal/server"
)

func main() {
	address := os.Getenv("GRPC_ADDRESS")
	if address == "" {
		address = ":50051"
	}

	srv, err := server.New(address)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	go func() {
		log.Printf("gRPC server listening on %s", srv.Address())
		if err := srv.Serve(); err != nil {
			log.Fatalf("server stopped: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down")
	srv.Stop()
}
