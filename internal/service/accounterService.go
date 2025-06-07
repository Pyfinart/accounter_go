package service

import (
	"context"

	v1 "accounter_go/api/accounter/v1"
	"accounter_go/internal/biz"
)

// GreeterService is a greeter service.
type AccounterService struct {
	v1.UnimplementedAccounterServer

	uc *biz.AccounterUseCase
}

// NewGreeterService new a greeter service.
func NewAccounterService(uc *biz.AccounterUseCase) *AccounterService {
	return &AccounterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *AccounterService) Add(ctx context.Context, in *v1.AddRequest) (*v1.AddReply, error) {
	_, err := s.uc.CreateAccounter(ctx, &biz.Accounter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.AddReply{}, nil
}
