package server

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	flashcardv1 "github.com/Shivam583-hue/trueflashcard/server/gen/flashcard/v1"
	"github.com/Shivam583-hue/trueflashcard/server/internal/auth"
	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
	"github.com/Shivam583-hue/trueflashcard/server/internal/service"
)

type Server struct {
	grpc     *grpc.Server
	listener net.Listener
}

func New(address string, q dbgen.Querier, tx service.Transactor, sessions *auth.SessionManager) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	var opts []grpc.ServerOption
	if sessions != nil {
		opts = append(opts, grpc.UnaryInterceptor(auth.UnaryAuthInterceptor(sessions)))
	}
	grpcServer := grpc.NewServer(opts...)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	if q != nil {
		flashcardv1.RegisterFolderServiceServer(grpcServer, service.NewFolderService(q))
		flashcardv1.RegisterDeckServiceServer(grpcServer, service.NewDeckService(q, tx))
		flashcardv1.RegisterFlashcardServiceServer(grpcServer, service.NewFlashcardService(q, tx))
	}

	reflection.Register(grpcServer)

	return &Server{grpc: grpcServer, listener: listener}, nil
}

func (s *Server) Address() string {
	return s.listener.Addr().String()
}

func (s *Server) Serve() error {
	return s.grpc.Serve(s.listener)
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
