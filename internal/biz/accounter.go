package biz

import (
	"context"
	"time"

	v1 "accounter_go/api/accounter/v1"
	"github.com/go-kratos/kratos/v2/log"
)

// Accounter is a Accounter model.
type Accounter struct {
	TransactionID int64
	UserID        int64
	Type          v1.Type
	Category      v1.Category
	Desc          string
	Amount        float64
	Date          time.Time
}

// AccounterRepo is a Accounter repo.
type AccounterRepo interface {
	Save(context.Context, *Accounter) (*Accounter, error)
	Update(context.Context, *Accounter) (*Accounter, error)
	FindByID(context.Context, int64) (*Accounter, error)
	ListByUserID(context.Context, int64) ([]*Accounter, error)
	ListAll(context.Context) ([]*Accounter, error)
}

// AccounterUseCase is a Accounter usecase.
type AccounterUseCase struct {
	repo AccounterRepo
	log  *log.Helper
}

// NewAccounterUsecase new a Accounter usecase.
func NewAccounterUsecase(repo AccounterRepo, logger log.Logger) *AccounterUseCase {
	return &AccounterUseCase{repo: repo, log: log.NewHelper(logger)}
}

// CreateAccounter creates a Accounter, and returns the new Accounter.
func (uc *AccounterUseCase) CreateAccounter(ctx context.Context, g *Accounter) (*Accounter, error) {
	uc.log.WithContext(ctx).Infof("CreateAccounter: %v", g.Desc)
	return uc.repo.Save(ctx, g)
}
