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
	ListWithFilters(context.Context, *ListFilter) ([]*Accounter, int32, error)
	Delete(context.Context, int64) error
	GetStats(context.Context, *StatsFilter) (*Stats, error)
}

// AccounterUseCase is a Accounter usecase.
type AccounterUseCase struct {
	repo AccounterRepo
	Log  *log.Helper
}

// NewAccounterUsecase new a Accounter usecase.
func NewAccounterUsecase(repo AccounterRepo, logger log.Logger) *AccounterUseCase {
	return &AccounterUseCase{repo: repo, Log: log.NewHelper(logger)}
}

// ListFilter represents filters for listing transactions
type ListFilter struct {
	UserID    int64
	Type      *v1.Type
	Category  *v1.Category
	StartDate *time.Time
	EndDate   *time.Time
	Page      int32
	PageSize  int32
}

// StatsFilter represents filters for getting statistics
type StatsFilter struct {
	UserID    int64
	StartDate *time.Time
	EndDate   *time.Time
}

// CategoryStat represents statistics for a category
type CategoryStat struct {
	Category     v1.Category
	CategoryName string
	Amount       float64
	Count        int32
}

// Stats represents financial statistics
type Stats struct {
	TotalIncome       float64
	TotalExpense      float64
	Balance           float64
	IncomeByCategory  []*CategoryStat
	ExpenseByCategory []*CategoryStat
}

// CreateAccounter creates a Accounter, and returns the new Accounter.
func (uc *AccounterUseCase) CreateAccounter(ctx context.Context, g *Accounter) (*Accounter, error) {
	uc.Log.WithContext(ctx).Infof("CreateAccounter: %v", g.Desc)
	return uc.repo.Save(ctx, g)
}

// ListAccounters lists accounters with filters
func (uc *AccounterUseCase) ListAccounters(ctx context.Context, filter *ListFilter) ([]*Accounter, int32, error) {
	uc.Log.WithContext(ctx).Infof("ListAccounters with filters")
	return uc.repo.ListWithFilters(ctx, filter)
}

// DeleteAccounter deletes an accounter by ID
func (uc *AccounterUseCase) DeleteAccounter(ctx context.Context, id int64) error {
	uc.Log.WithContext(ctx).Infof("DeleteAccounter: %d", id)
	return uc.repo.Delete(ctx, id)
}

// GetStats gets financial statistics
func (uc *AccounterUseCase) GetStats(ctx context.Context, filter *StatsFilter) (*Stats, error) {
	uc.Log.WithContext(ctx).Infof("GetStats")
	return uc.repo.GetStats(ctx, filter)
}
