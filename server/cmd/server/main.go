package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shivam583-hue/trueflashcard/server/internal/auth"
	"github.com/Shivam583-hue/trueflashcard/server/internal/connectapi"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
	"github.com/Shivam583-hue/trueflashcard/server/internal/httpauth"
	"github.com/Shivam583-hue/trueflashcard/server/internal/server"
	"github.com/Shivam583-hue/trueflashcard/server/internal/service"
)

func main() {
	grpcAddress := envOr("GRPC_ADDRESS", ":50051")
	httpAddress := envOr("HTTP_ADDRESS", ":8080")

	ctx := context.Background()

	var querier dbgen.Querier
	var tx service.Transactor
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
		tx = store
		log.Println("database connection established")
	}

	sessions := buildSessionManager()
	httpServer := buildHTTPServer(httpAddress, querier, tx, sessions)
	if httpServer != nil {
		go func() {
			log.Printf("HTTP server listening on %s (Connect API + auth)", httpAddress)
			if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("HTTP server stopped: %v", err)
			}
		}()
	}

	srv, err := server.New(grpcAddress, querier, tx, sessions)
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
	if httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
	}
	srv.Stop()
}

func buildSessionManager() *auth.SessionManager {
	sessions, err := auth.NewSessionManager()
	if errors.Is(err, auth.ErrMissingJWTSecret) {
		log.Println("JWT_SECRET not set; sessions and authenticated RPCs are disabled")
		return nil
	}
	if err != nil {
		log.Fatalf("failed to initialize sessions: %v", err)
	}
	return sessions
}

func buildHTTPServer(address string, querier dbgen.Querier, tx service.Transactor, sessions *auth.SessionManager) *http.Server {
	if querier == nil || sessions == nil {
		return nil
	}

	appURL := envOr("APP_URL", "http://localhost:3000")
	mux := http.NewServeMux()

	mux.Handle("/", connectapi.NewHandler(querier, tx, sessions, appURL))
	log.Println("Connect API enabled")

	oauth, err := auth.NewGoogleOAuth()
	switch {
	case errors.Is(err, auth.ErrMissingOAuthConfig):
		log.Println("Google OAuth env not set; login flow is disabled")
	case err != nil:
		log.Fatalf("failed to initialize Google OAuth: %v", err)
	default:
		mux.Handle("/auth/", httpauth.NewHandler(oauth, sessions, querier).Routes())
		log.Println("Google OAuth login enabled")
	}

	return &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
