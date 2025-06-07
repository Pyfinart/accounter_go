package service

import (
	"context"
	"time"

	v1 "accounter_go/api/accounter/v1"
	"accounter_go/internal/biz"
)

// AccounterService is a accounter service.
type AccounterService struct {
	v1.UnimplementedAccounterServer

	uc *biz.AccounterUseCase
}

// NewAccounterService new a accounter service.
func NewAccounterService(uc *biz.AccounterUseCase) *AccounterService {
	return &AccounterService{uc: uc}
}

// Add implements accounter.AccounterServer.
func (s *AccounterService) Add(ctx context.Context, in *v1.AddRequest) (*v1.AddReply, error) {
	// Parse date string to time.Time
	transactionDate, err := time.Parse("2006-01-02", in.Date)
	if err != nil {
		// If date parsing fails, use current time
		transactionDate = time.Now()
	}

	// Create biz.Accounter from request
	accounter := &biz.Accounter{
		UserID:   1, // TODO: Get from context/auth
		Type:     in.Type,
		Category: in.Category,
		Desc:     in.Desc,
		Amount:   in.Amount,
		Date:     transactionDate,
	}

	_, err = s.uc.CreateAccounter(ctx, accounter)
	if err != nil {
		return nil, err
	}
	
	return &v1.AddReply{}, nil
}
