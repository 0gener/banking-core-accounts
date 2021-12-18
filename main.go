package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"

	"github.com/0gener/banking-core-accounts/data"
	pb "github.com/0gener/banking-core-accounts/proto"
	"google.golang.org/grpc"
)

type server struct {
	repo data.Repository
	pb.UnimplementedAccountsServiceServer
}

var ErrInvalidUserId = errors.New("invalid user_id")
var ErrInvalidCurrency = errors.New("invalid currency")

func (s *server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	if req.UserId == "" {
		return nil, ErrInvalidUserId
	}
	if req.Currency == "" {
		return nil, ErrInvalidCurrency
	}

	account, _ := s.repo.Save(data.AccountEntity{
		UserId:        req.UserId,
		AccountNumber: strconv.FormatUint(rand.Uint64(), 10),
		Currency:      req.Currency,
	})

	return &pb.CreateAccountResponse{
		Account: &pb.Account{
			AccountNumber: account.AccountNumber,
			Currency:      account.Currency,
		},
	}, nil
}

func (s *server) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	if req.UserId == "" {
		return nil, ErrInvalidUserId
	}

	account, _ := s.repo.FindByUserId(req.UserId)

	if account == nil {
		return &pb.GetAccountResponse{}, nil
	}

	return &pb.GetAccountResponse{
		Account: &pb.Account{
			AccountNumber: account.AccountNumber,
			Currency:      account.Currency,
		},
	}, nil
}

func main() {
	s := &server{
		repo: data.NewInMemoryRepository(),
	}
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
