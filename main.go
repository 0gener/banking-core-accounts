package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/0gener/banking-core-accounts/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAccountsServiceServer
}

func (s *server) CreateAccount(context.Context, *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	return nil, nil
}

func main() {
	s := &server{}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 5000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAccountsServiceServer(grpcServer, s)

	log.Print("Server listening on http://localhost:5000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
