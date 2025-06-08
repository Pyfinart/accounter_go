// This file contains the database-based implementation of AccounterRepo
// It's kept for future use when switching back to database storage
// To use this implementation:
// 1. Uncomment NewAccounterDbRepo in data.go ProviderSet
// 2. Comment out NewAccounterFileRepo in data.go ProviderSet
// 3. Ensure database configuration is properly set up

package data

import (
	"context"
	"fmt"

	v1 "accounter_go/api/accounter/v1"
	"accounter_go/internal/biz"
	"accounter_go/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type accounterDbRepo struct {
	data *Data
	log  *log.Helper
}

// NewAccounterDbRepo creates a new database-based AccounterRepo
// This is kept for future use when switching back to database storage
func NewAccounterDbRepo(data *Data, logger log.Logger) biz.AccounterRepo {
	return &accounterDbRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *accounterDbRepo) Save(ctx context.Context, accounter *biz.Accounter) (*biz.Accounter, error) {
	// Convert biz.Accounter to model.AccounterTransaction
	transaction := &model.AccounterTransaction{
		UserID:          accounter.UserID,
		CategoryID:      int(accounter.Category),
		CurrencyID:      1, // Default to CNY
		TransactionType: int8(accounter.Type),
		Amount:          accounter.Amount,
		TransactionDate: accounter.Date,
		Note:            &accounter.Desc,
	}

	if err := r.data.db.WithContext(ctx).Create(transaction).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to save accounter: %v", err)
		return nil, err
	}

	// Convert back to biz.Accounter
	result := &biz.Accounter{
		TransactionID: transaction.TransactionID,
		UserID:        transaction.UserID,
		Type:          v1.Type(transaction.TransactionType),
		Category:      v1.Category(transaction.CategoryID),
		Desc:          *transaction.Note,
		Amount:        transaction.Amount,
		Date:          transaction.TransactionDate,
	}

	return result, nil
}

func (r *accounterDbRepo) Update(ctx context.Context, accounter *biz.Accounter) (*biz.Accounter, error) {
	transaction := &model.AccounterTransaction{
		TransactionID:   accounter.TransactionID,
		UserID:          accounter.UserID,
		CategoryID:      int(accounter.Category),
		CurrencyID:      1, // Default to CNY
		TransactionType: int8(accounter.Type),
		Amount:          accounter.Amount,
		TransactionDate: accounter.Date,
		Note:            &accounter.Desc,
	}

	if err := r.data.db.WithContext(ctx).Save(transaction).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update accounter: %v", err)
		return nil, err
	}

	return accounter, nil
}

func (r *accounterDbRepo) FindByID(ctx context.Context, id int64) (*biz.Accounter, error) {
	var transaction model.AccounterTransaction
	if err := r.data.db.WithContext(ctx).First(&transaction, id).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to find accounter by id %d: %v", id, err)
		return nil, err
	}

	result := &biz.Accounter{
		TransactionID: transaction.TransactionID,
		UserID:        transaction.UserID,
		Type:          v1.Type(transaction.TransactionType),
		Category:      v1.Category(transaction.CategoryID),
		Desc:          *transaction.Note,
		Amount:        transaction.Amount,
		Date:          transaction.TransactionDate,
	}

	return result, nil
}

func (r *accounterDbRepo) ListByUserID(ctx context.Context, userID int64) ([]*biz.Accounter, error) {
	var transactions []model.AccounterTransaction
	if err := r.data.db.WithContext(ctx).Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list accounters by user id %d: %v", userID, err)
		return nil, err
	}

	var results []*biz.Accounter
	for _, transaction := range transactions {
		note := ""
		if transaction.Note != nil {
			note = *transaction.Note
		}
		result := &biz.Accounter{
			TransactionID: transaction.TransactionID,
			UserID:        transaction.UserID,
			Type:          v1.Type(transaction.TransactionType),
			Category:      v1.Category(transaction.CategoryID),
			Desc:          note,
			Amount:        transaction.Amount,
			Date:          transaction.TransactionDate,
		}
		results = append(results, result)
	}

	return results, nil
}

func (r *accounterDbRepo) ListAll(ctx context.Context) ([]*biz.Accounter, error) {
	var transactions []model.AccounterTransaction
	if err := r.data.db.WithContext(ctx).Find(&transactions).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list all accounters: %v", err)
		return nil, err
	}

	var results []*biz.Accounter
	for _, transaction := range transactions {
		note := ""
		if transaction.Note != nil {
			note = *transaction.Note
		}
		result := &biz.Accounter{
			TransactionID: transaction.TransactionID,
			UserID:        transaction.UserID,
			Type:          v1.Type(transaction.TransactionType),
			Category:      v1.Category(transaction.CategoryID),
			Desc:          note,
			Amount:        transaction.Amount,
			Date:          transaction.TransactionDate,
		}
		results = append(results, result)
	}

	return results, nil
}

// ListWithFilters - Database implementation (placeholder for future use)
func (r *accounterDbRepo) ListWithFilters(ctx context.Context, filter *biz.ListFilter) ([]*biz.Accounter, int32, error) {
	// TODO: Implement database version with proper SQL queries
	return nil, 0, fmt.Errorf("database implementation not yet available")
}

// Delete - Database implementation (placeholder for future use)
func (r *accounterDbRepo) Delete(ctx context.Context, id int64) error {
	// TODO: Implement database version
	return fmt.Errorf("database implementation not yet available")
}

// GetStats - Database implementation (placeholder for future use)
func (r *accounterDbRepo) GetStats(ctx context.Context, filter *biz.StatsFilter) (*biz.Stats, error) {
	// TODO: Implement database version with proper SQL queries
	return nil, fmt.Errorf("database implementation not yet available")
}
