package server

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpc     *grpc.Server
	listener net.Listener
}

func New(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

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
