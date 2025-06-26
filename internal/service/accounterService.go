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

	result, err := s.uc.CreateAccounter(ctx, accounter)
	if err != nil {
		return nil, err
	}

	return &v1.AddReply{
		Id:      result.TransactionID,
		Message: "Transaction created successfully",
	}, nil
}

// List implements accounter.AccounterServer.
func (s *AccounterService) List(ctx context.Context, in *v1.ListRequest) (*v1.ListReply, error) {
	s.uc.Log.Errorf("ListAccounters with filters, params: %+v", in)

	const defaultPageSize = 20
	const defaultPage = 1

	filter := &biz.ListFilter{
		UserID:   1, // TODO: Get from context/auth
		Page:     in.Page,
		PageSize: in.PageSize,
	}

	if in.Type != v1.Type_None {
		filter.Type = &in.Type
	}
	if in.Category != v1.Category_Default {
		filter.Category = &in.Category
	}

	// Parse date filters
	if in.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", in.StartDate); err == nil {
			filter.StartDate = &startDate
		}
	}
	if in.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", in.EndDate); err == nil {
			filter.EndDate = &endDate
		}
	}

	// Set default pagination
	if filter.Page <= 0 {
		filter.Page = defaultPage
	}
	if filter.PageSize <= 0 {
		filter.PageSize = defaultPageSize
	}

	accounters, total, err := s.uc.ListAccounters(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	transactions := make([]*v1.Transaction, len(accounters))
	for i, acc := range accounters {
		transactions[i] = &v1.Transaction{
			Id:        acc.TransactionID,
			Type:      acc.Type,
			Category:  acc.Category,
			Desc:      acc.Desc,
			Amount:    acc.Amount,
			Date:      acc.Date.Format("2006-01-02"),
			CreatedAt: acc.Date.Format("2006-01-02 15:04:05"),
		}
	}

	return &v1.ListReply{
		Transactions: transactions,
		Total:        total,
		Page:         filter.Page,
		PageSize:     filter.PageSize,
	}, nil
}

// Stats implements accounter.AccounterServer.
func (s *AccounterService) Stats(ctx context.Context, in *v1.StatsRequest) (*v1.StatsReply, error) {
	filter := &biz.StatsFilter{
		UserID: 1, // TODO: Get from context/auth
	}

	// Parse date filters
	if in.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", in.StartDate); err == nil {
			filter.StartDate = &startDate
		}
	}
	if in.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", in.EndDate); err == nil {
			filter.EndDate = &endDate
		}
	}

	stats, err := s.uc.GetStats(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	incomeByCategory := make([]*v1.CategoryStats, len(stats.IncomeByCategory))
	for i, cat := range stats.IncomeByCategory {
		incomeByCategory[i] = &v1.CategoryStats{
			Category:     cat.Category,
			CategoryName: cat.CategoryName,
			Amount:       cat.Amount,
			Count:        cat.Count,
		}
	}

	expenseByCategory := make([]*v1.CategoryStats, len(stats.ExpenseByCategory))
	for i, cat := range stats.ExpenseByCategory {
		expenseByCategory[i] = &v1.CategoryStats{
			Category:     cat.Category,
			CategoryName: cat.CategoryName,
			Amount:       cat.Amount,
			Count:        cat.Count,
		}
	}

	return &v1.StatsReply{
		TotalIncome:       stats.TotalIncome,
		TotalExpense:      stats.TotalExpense,
		Balance:           stats.Balance,
		IncomeByCategory:  incomeByCategory,
		ExpenseByCategory: expenseByCategory,
	}, nil
}

// Delete implements accounter.AccounterServer.
func (s *AccounterService) Delete(ctx context.Context, in *v1.DeleteRequest) (*v1.DeleteReply, error) {
	err := s.uc.DeleteAccounter(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteReply{
		Message: "Transaction deleted successfully",
	}, nil
}
