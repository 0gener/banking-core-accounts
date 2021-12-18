package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "github.com/0gener/banking-core-accounts/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAccountsServiceServer
}

var ErrInvalidUserId = errors.New("invalid user_id")
var ErrInvalidCurrency = errors.New("invalid currency")

var accounts map[string]*pb.Account

func (s *server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	if req.UserId == "" {
		return nil, ErrInvalidUserId
	}
	if req.Currency == "" {
		return nil, ErrInvalidCurrency
	}

	if accounts == nil {
		accounts = make(map[string]*pb.Account)
	}

	accounts[req.UserId] = &pb.Account{
		AccountNumber: strconv.FormatUint(rand.Uint64(), 10),
		Currency:      req.Currency,
	}

	return &pb.CreateAccountResponse{
		Account: accounts[req.UserId],
	}, nil
}

func (s *server) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	if req.UserId == "" {
		return nil, ErrInvalidUserId
	}

	return &pb.GetAccountResponse{
		Account: accounts[req.UserId],
	}, nil
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
