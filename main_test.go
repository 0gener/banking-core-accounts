package main

import (
	"context"
	"testing"

	pb "github.com/0gener/banking-core-accounts/proto"
)

func TestCreateAccount(t *testing.T) {
	testCases := []struct {
		name             string
		request          *pb.CreateAccountRequest
		expectedResponse *pb.CreateAccountResponse
		expectedError    error
	}{
		{
			name: "valid request",
			request: &pb.CreateAccountRequest{
				UserId:   "1234",
				Currency: "EUR",
			},
			expectedResponse: &pb.CreateAccountResponse{
				Account: &pb.Account{
					Currency: "EUR",
				},
			},
			expectedError: nil,
		},
		{
			name: "request without user_id",
			request: &pb.CreateAccountRequest{
				UserId:   "",
				Currency: "EUR",
			},
			expectedResponse: nil,
			expectedError:    ErrInvalidUserId,
		},
		{
			name: "request without currency",
			request: &pb.CreateAccountRequest{
				UserId:   "1234",
				Currency: "",
			},
			expectedResponse: nil,
			expectedError:    ErrInvalidCurrency,
		},
	}

	s := server{}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, err := s.CreateAccount(context.Background(), tc.request)
			if tc.expectedError != nil {
				if err != tc.expectedError {
					t.Errorf("expected error: '%v', got '%v'", tc.expectedError, err)
				}
				if res != nil {
					t.Errorf("expected no response when error occured, got %v", res)
				}
			} else {
				if res.Account.AccountNumber == "" {
					t.Errorf("expected a non-blank AccountNumber")
				}
				if res.Account.Currency != tc.request.Currency {
					t.Errorf("expected Currency '%s', got '%s'", tc.request.Currency, res.Account.Currency)
				}
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	testCases := []struct {
		name             string
		request          *pb.GetAccountRequest
		expectedResponse *pb.GetAccountResponse
		expectedError    error
	}{
		{
			name: "valid request",
			request: &pb.GetAccountRequest{
				UserId: "1234",
			},
			expectedResponse: &pb.GetAccountResponse{
				Account: &pb.Account{
					Currency: "EUR",
				},
			},
			expectedError: nil,
		},
		{
			name: "request without user_id",
			request: &pb.GetAccountRequest{
				UserId: "",
			},
			expectedResponse: nil,
			expectedError:    ErrInvalidUserId,
		},
	}

	s := server{}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res, err := s.GetAccount(context.Background(), tc.request)
			if tc.expectedError != nil {
				if err != tc.expectedError {
					t.Errorf("expected error: '%v', got '%v'", tc.expectedError, err)
				}
				if res != nil {
					t.Errorf("expected no response when error occured, got %v", res)
				}
			} else {
				if res.Account.AccountNumber == "" {
					t.Errorf("expected a non-blank AccountNumber")
				}
			}
		})
	}
}
