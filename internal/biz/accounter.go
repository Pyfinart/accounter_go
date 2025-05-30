package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Greeter is a Greeter model.
type Accounter struct {
	Hello string
}

// GreeterRepo is a Greater repo.
type AccounterRepo interface {
	Save(context.Context, *Accounter) (*Accounter, error)
	Update(context.Context, *Accounter) (*Accounter, error)
	FindByID(context.Context, int64) (*Accounter, error)
	ListByHello(context.Context, string) ([]*Accounter, error)
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
	uc.log.WithContext(ctx).Infof("CreateAccounter: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}
