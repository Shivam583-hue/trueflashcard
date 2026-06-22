package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shivam583-hue/trueflashcard/server/internal/db"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
	"github.com/Shivam583-hue/trueflashcard/server/internal/server"
)

func main() {
	address := os.Getenv("GRPC_ADDRESS")
	if address == "" {
		address = ":50051"
	}

	ctx := context.Background()

	var querier dbgen.Querier
	store, err := db.Connect(ctx)
	switch {
	case errors.Is(err, db.ErrMissingDatabaseURL):
		log.Println("DATABASE_URL not set; starting with health check only")
	case err != nil:
		log.Fatalf("failed to connect to database: %v", err)
	default:
		defer store.Close()
		if err := store.VerifyConnectivity(ctx); err != nil {
			log.Fatalf("database connectivity check failed: %v", err)
		}
		querier = store
		log.Println("database connection established")
	}

	srv, err := server.New(address, querier)
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
