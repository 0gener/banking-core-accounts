package main

import (
	"context"
	"errors"
	"testing"

	"github.com/0gener/banking-core-accounts/data"
	pb "github.com/0gener/banking-core-accounts/proto"
)

type createAccountTestCase struct {
	name               string
	request            *pb.CreateAccountRequest
	expectedResponse   *pb.CreateAccountResponse
	expectedError      error
	repositoryResponse *data.AccountEntity
	repositoryError    error
}

type getAccountTestCase struct {
	name               string
	request            *pb.GetAccountRequest
	expectedResponse   *pb.GetAccountResponse
	expectedError      error
	repositoryResponse *data.AccountEntity
	repositoryError    error
}

type testRepository struct {
	t                     *testing.T
	createAccountTestCase createAccountTestCase
	getAccountTestCase    getAccountTestCase
}

func (repo *testRepository) Save(entity data.AccountEntity) (*data.AccountEntity, error) {
	if entity.UserId != repo.createAccountTestCase.request.UserId {
		repo.t.Errorf("expected userId %s, got %s", repo.createAccountTestCase.request.UserId, entity.UserId)
	}
	if entity.Currency != repo.createAccountTestCase.request.Currency {
		repo.t.Errorf("expected currency %s, got %s", repo.createAccountTestCase.request.Currency, entity.Currency)
	}
	return repo.createAccountTestCase.repositoryResponse, repo.createAccountTestCase.repositoryError
}

func (repo *testRepository) FindByUserId(userId string) (*data.AccountEntity, error) {
	return repo.getAccountTestCase.repositoryResponse, repo.getAccountTestCase.repositoryError
}

func TestCreateAccount(t *testing.T) {
	testCases := []createAccountTestCase{
		{
			name: "valid request",
			request: &pb.CreateAccountRequest{
				UserId:   "1234",
				Currency: "EUR",
			},
			expectedResponse: &pb.CreateAccountResponse{
				Account: &pb.Account{
					AccountNumber: "123456",
					Currency:      "EUR",
				},
			},
			expectedError: nil,
			repositoryResponse: &data.AccountEntity{
				UserId:        "1234",
				AccountNumber: "123456",
				Currency:      "EUR",
			},
			repositoryError: nil,
		},
		{
			name: "request without user_id",
			request: &pb.CreateAccountRequest{
				UserId:   "",
				Currency: "EUR",
			},
			expectedResponse:   nil,
			expectedError:      ErrInvalidUserId,
			repositoryResponse: nil,
			repositoryError:    nil,
		},
		{
			name: "request without currency",
			request: &pb.CreateAccountRequest{
				UserId:   "1234",
				Currency: "",
			},
			expectedResponse:   nil,
			expectedError:      ErrInvalidCurrency,
			repositoryResponse: nil,
			repositoryError:    nil,
		},
		{
			name: "repository returns an error",
			request: &pb.CreateAccountRequest{
				UserId:   "1234",
				Currency: "EUR",
			},
			expectedResponse:   nil,
			expectedError:      ErrRepository,
			repositoryResponse: nil,
			repositoryError:    errors.New("just a random error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := server{
				repo: &testRepository{
					t:                     t,
					createAccountTestCase: tc,
				},
			}

			res, err := s.CreateAccount(context.Background(), tc.request)
			if tc.expectedError != nil {
				if err != tc.expectedError {
					t.Errorf("expected error: '%v', got '%v'", tc.expectedError, err)
				}
				if res != nil {
					t.Errorf("expected no response when error occured, got %v", res)
				}
			} else {
				if res == tc.expectedResponse {
					t.Errorf("expected the response to be %v, got %v", tc.expectedResponse, res)
				}
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	testCases := []getAccountTestCase{
		{
			name: "valid request",
			request: &pb.GetAccountRequest{
				UserId: "1234",
			},
			expectedResponse: &pb.GetAccountResponse{
				Account: &pb.Account{
					AccountNumber: "123456",
					Currency:      "EUR",
				},
			},
			expectedError: nil,
			repositoryResponse: &data.AccountEntity{
				UserId:        "1234",
				AccountNumber: "123456",
				Currency:      "EUR",
			},
			repositoryError: nil,
		},
		{
			name: "request without user_id",
			request: &pb.GetAccountRequest{
				UserId: "",
			},
			expectedResponse:   nil,
			expectedError:      ErrInvalidUserId,
			repositoryResponse: nil,
			repositoryError:    nil,
		},
		{
			name: "repository returns no account",
			request: &pb.GetAccountRequest{
				UserId: "1234",
			},
			expectedResponse:   &pb.GetAccountResponse{},
			expectedError:      nil,
			repositoryResponse: nil,
			repositoryError:    nil,
		},
		{
			name: "repository returns an error",
			request: &pb.GetAccountRequest{
				UserId: "1234",
			},
			expectedResponse:   nil,
			expectedError:      ErrRepository,
			repositoryResponse: nil,
			repositoryError:    errors.New("just a random error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := server{
				repo: &testRepository{
					t:                  t,
					getAccountTestCase: tc,
				},
			}

			res, err := s.GetAccount(context.Background(), tc.request)
			if tc.expectedError != nil {
				if err != tc.expectedError {
					t.Errorf("expected error: '%v', got '%v'", tc.expectedError, err)
				}
				if res != nil {
					t.Errorf("expected no response when error occured, got %v", res)
				}
			} else {
				if res == tc.expectedResponse {
					t.Errorf("expected the response to be %v, got %v", tc.expectedResponse, res)
				}
			}
		})
	}
}
