package grpcserver

import (
	"log"
	"net"

	pb "go-fruit-cart/newproto"

	"google.golang.org/grpc"
)

// GRPCServer implements gRPC server for rating management service
type GRPCServer struct {
	port       string
	grpcServer pb.FruitCartManagementServiceServer
}

// NewGRPCServer returns new instance of rms gRPC server
func NewGRPCServer(port string, grpcServer pb.FruitCartManagementServiceServer) *GRPCServer {
	return &GRPCServer{
		port:       port,
		grpcServer: grpcServer,
	}
}

// ListenAndServe servers the gRPC server on specified port
func (s *GRPCServer) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFruitCartManagementServiceServer(grpcServer, s.grpcServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Println("error when serving grpc server: ", s.port, err)
		return err
	}

	return nil
}
