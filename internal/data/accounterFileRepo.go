package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	v1 "accounter_go/api/accounter/v1"
	"accounter_go/internal/biz"
	"accounter_go/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// FileAccounterData represents the structure stored in JSON file
type FileAccounterData struct {
	TransactionID int64     `json:"transaction_id"`
	UserID        int64     `json:"user_id"`
	Type          int32     `json:"type"`
	Category      int32     `json:"category"`
	Desc          string    `json:"desc"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
}

// FileAccounterStorage manages the file storage operations
type FileAccounterStorage struct {
	filePath string
	data     []FileAccounterData
	mutex    sync.RWMutex
	nextID   int64
	log      *log.Helper
}

type accounterFileRepo struct {
	storage *FileAccounterStorage
	log     *log.Helper
}

// NewAccounterFileRepo creates a new file-based AccounterRepo
func NewAccounterFileRepo(c *conf.Data, logger log.Logger) biz.AccounterRepo {
	// Get file storage config or use defaults
	dataDir := "./data"
	fileName := "accounters.json"
	
	if c.FileStorage != nil {
		if c.FileStorage.DataDir != "" {
			dataDir = c.FileStorage.DataDir
		}
		if c.FileStorage.AccounterFile != "" {
			fileName = c.FileStorage.AccounterFile
		}
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.NewHelper(logger).Errorf("Failed to create data directory: %v", err)
	}

	storage := &FileAccounterStorage{
		filePath: filepath.Join(dataDir, fileName),
		data:     make([]FileAccounterData, 0),
		nextID:   1,
		log:      log.NewHelper(logger),
	}

	// Load existing data
	storage.loadFromFile()

	log.NewHelper(logger).Infof("Initialized file storage at: %s", storage.filePath)

	return &accounterFileRepo{
		storage: storage,
		log:     log.NewHelper(logger),
	}
}

func (s *FileAccounterStorage) loadFromFile() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// File doesn't exist, start with empty data
		return
	}

	content, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		s.log.Errorf("Failed to read file %s: %v", s.filePath, err)
		return
	}

	if len(content) == 0 {
		return
	}

	if err := json.Unmarshal(content, &s.data); err != nil {
		s.log.Errorf("Failed to unmarshal data from file %s: %v", s.filePath, err)
		return
	}

	// Find the next ID
	for _, item := range s.data {
		if item.TransactionID >= s.nextID {
			s.nextID = item.TransactionID + 1
		}
	}

	s.log.Infof("Loaded %d records from file, next ID: %d", len(s.data), s.nextID)
}

func (s *FileAccounterStorage) saveToFile() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	content, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := os.WriteFile(s.filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %v", s.filePath, err)
	}

	return nil
}

func (r *accounterFileRepo) Save(ctx context.Context, accounter *biz.Accounter) (*biz.Accounter, error) {
	r.storage.mutex.Lock()

	// Generate new ID
	newID := r.storage.nextID
	r.storage.nextID++

	// Create file data structure
	fileData := FileAccounterData{
		TransactionID: newID,
		UserID:        accounter.UserID,
		Type:          int32(accounter.Type),
		Category:      int32(accounter.Category),
		Desc:          accounter.Desc,
		Amount:        accounter.Amount,
		Date:          accounter.Date,
		CreatedAt:     time.Now(),
	}

	// Add to in-memory data
	r.storage.data = append(r.storage.data, fileData)
	r.storage.mutex.Unlock()

	// Save to file
	if err := r.storage.saveToFile(); err != nil {
		r.log.WithContext(ctx).Errorf("Failed to save to file: %v", err)
		return nil, err
	}

	// Return the saved accounter with new ID
	result := &biz.Accounter{
		TransactionID: newID,
		UserID:        accounter.UserID,
		Type:          accounter.Type,
		Category:      accounter.Category,
		Desc:          accounter.Desc,
		Amount:        accounter.Amount,
		Date:          accounter.Date,
	}

	r.log.WithContext(ctx).Infof("Saved accounter with ID: %d", newID)
	return result, nil
}

func (r *accounterFileRepo) Update(ctx context.Context, accounter *biz.Accounter) (*biz.Accounter, error) {
	r.storage.mutex.Lock()
	defer r.storage.mutex.Unlock()

	// Find and update the record
	for i, item := range r.storage.data {
		if item.TransactionID == accounter.TransactionID {
			r.storage.data[i] = FileAccounterData{
				TransactionID: accounter.TransactionID,
				UserID:        accounter.UserID,
				Type:          int32(accounter.Type),
				Category:      int32(accounter.Category),
				Desc:          accounter.Desc,
				Amount:        accounter.Amount,
				Date:          accounter.Date,
				CreatedAt:     item.CreatedAt, // Keep original creation time
			}

			// Save to file
			if err := r.storage.saveToFile(); err != nil {
				r.log.WithContext(ctx).Errorf("Failed to save to file after update: %v", err)
				return nil, err
			}

			r.log.WithContext(ctx).Infof("Updated accounter with ID: %d", accounter.TransactionID)
			return accounter, nil
		}
	}

	return nil, fmt.Errorf("accounter with ID %d not found", accounter.TransactionID)
}

func (r *accounterFileRepo) FindByID(ctx context.Context, id int64) (*biz.Accounter, error) {
	r.storage.mutex.RLock()
	defer r.storage.mutex.RUnlock()

	for _, item := range r.storage.data {
		if item.TransactionID == id {
			result := &biz.Accounter{
				TransactionID: item.TransactionID,
				UserID:        item.UserID,
				Type:          v1.Type(item.Type),
				Category:      v1.Category(item.Category),
				Desc:          item.Desc,
				Amount:        item.Amount,
				Date:          item.Date,
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("accounter with ID %d not found", id)
}

func (r *accounterFileRepo) ListByUserID(ctx context.Context, userID int64) ([]*biz.Accounter, error) {
	r.storage.mutex.RLock()
	defer r.storage.mutex.RUnlock()

	var results []*biz.Accounter
	for _, item := range r.storage.data {
		if item.UserID == userID {
			result := &biz.Accounter{
				TransactionID: item.TransactionID,
				UserID:        item.UserID,
				Type:          v1.Type(item.Type),
				Category:      v1.Category(item.Category),
				Desc:          item.Desc,
				Amount:        item.Amount,
				Date:          item.Date,
			}
			results = append(results, result)
		}
	}

	return results, nil
}

func (r *accounterFileRepo) ListAll(ctx context.Context) ([]*biz.Accounter, error) {
	r.storage.mutex.RLock()
	defer r.storage.mutex.RUnlock()

	var results []*biz.Accounter
	for _, item := range r.storage.data {
		result := &biz.Accounter{
			TransactionID: item.TransactionID,
			UserID:        item.UserID,
			Type:          v1.Type(item.Type),
			Category:      v1.Category(item.Category),
			Desc:          item.Desc,
			Amount:        item.Amount,
			Date:          item.Date,
		}
		results = append(results, result)
	}

	return results, nil
}

func (r *accounterFileRepo) ListWithFilters(ctx context.Context, filter *biz.ListFilter) ([]*biz.Accounter, int32, error) {
	r.storage.mutex.RLock()
	defer r.storage.mutex.RUnlock()

	var filtered []*biz.Accounter
	for _, item := range r.storage.data {
		// Apply filters
		if filter.UserID != 0 && item.UserID != filter.UserID {
			continue
		}
		if filter.Type != nil && int32(*filter.Type) != item.Type {
			continue
		}
		if filter.Category != nil && int32(*filter.Category) != item.Category {
			continue
		}
		if filter.StartDate != nil && item.Date.Before(*filter.StartDate) {
			continue
		}
		if filter.EndDate != nil && item.Date.After(*filter.EndDate) {
			continue
		}

		result := &biz.Accounter{
			TransactionID: item.TransactionID,
			UserID:        item.UserID,
			Type:          v1.Type(item.Type),
			Category:      v1.Category(item.Category),
			Desc:          item.Desc,
			Amount:        item.Amount,
			Date:          item.Date,
		}
		filtered = append(filtered, result)
	}

	total := int32(len(filtered))

	// Apply pagination
	start := (filter.Page - 1) * filter.PageSize
	end := start + filter.PageSize

	if start > total {
		return []*biz.Accounter{}, total, nil
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (r *accounterFileRepo) Delete(ctx context.Context, id int64) error {
	r.storage.mutex.Lock()
	defer r.storage.mutex.Unlock()

	// Find and remove the record
	for i, item := range r.storage.data {
		if item.TransactionID == id {
			r.storage.data = append(r.storage.data[:i], r.storage.data[i+1:]...)
			
			// Save to file
			if err := r.storage.saveToFile(); err != nil {
				r.log.WithContext(ctx).Errorf("Failed to save to file after delete: %v", err)
				return err
			}

			r.log.WithContext(ctx).Infof("Deleted accounter with ID: %d", id)
			return nil
		}
	}

	return fmt.Errorf("accounter with ID %d not found", id)
}

func (r *accounterFileRepo) GetStats(ctx context.Context, filter *biz.StatsFilter) (*biz.Stats, error) {
	r.storage.mutex.RLock()
	defer r.storage.mutex.RUnlock()

	var (
		totalIncome         float64
		totalExpense        float64
		incomeByCategory    = make(map[v1.Category]*biz.CategoryStat)
		expenseByCategory   = make(map[v1.Category]*biz.CategoryStat)
	)

	// Category names mapping
	categoryNames := map[v1.Category]string{
		v1.Category_Default:     "默认",
		v1.Category_Game:        "游戏",
		v1.Category_Food:        "餐饮",
		v1.Category_Travel:      "旅行",
		v1.Category_Education:   "教育",
		v1.Category_Health:      "健康",
		v1.Category_Shopping:    "购物",
		v1.Category_Other:       "其他",
		v1.Category_Transport:   "交通",
		v1.Category_Entertainment: "娱乐",
		v1.Category_Investment:  "投资",
		v1.Category_Loan:        "借款",
		v1.Category_Salary:      "工资",
		v1.Category_OtherIncome: "其他收入",
		v1.Category_App:         "应用",
		v1.Category_House:       "住房",
		v1.Category_Utility:     "水电费",
		v1.Category_Gift:        "礼物",
		v1.Category_Snacks:      "零食",
	}

	for _, item := range r.storage.data {
		// Apply filters
		if filter.UserID != 0 && item.UserID != filter.UserID {
			continue
		}
		if filter.StartDate != nil && item.Date.Before(*filter.StartDate) {
			continue
		}
		if filter.EndDate != nil && item.Date.After(*filter.EndDate) {
			continue
		}

		category := v1.Category(item.Category)
		categoryName := categoryNames[category]
		if categoryName == "" {
			categoryName = "未知"
		}

		if item.Type == int32(v1.Type_Income) {
			totalIncome += item.Amount
			if stat, exists := incomeByCategory[category]; exists {
				stat.Amount += item.Amount
				stat.Count++
			} else {
				incomeByCategory[category] = &biz.CategoryStat{
					Category:     category,
					CategoryName: categoryName,
					Amount:       item.Amount,
					Count:        1,
				}
			}
		} else if item.Type == int32(v1.Type_Expense) {
			totalExpense += item.Amount
			if stat, exists := expenseByCategory[category]; exists {
				stat.Amount += item.Amount
				stat.Count++
			} else {
				expenseByCategory[category] = &biz.CategoryStat{
					Category:     category,
					CategoryName: categoryName,
					Amount:       item.Amount,
					Count:        1,
				}
			}
		}
	}

	// Convert maps to slices
	var incomeStats []*biz.CategoryStat
	for _, stat := range incomeByCategory {
		incomeStats = append(incomeStats, stat)
	}

	var expenseStats []*biz.CategoryStat
	for _, stat := range expenseByCategory {
		expenseStats = append(expenseStats, stat)
	}

	return &biz.Stats{
		TotalIncome:       totalIncome,
		TotalExpense:      totalExpense,
		Balance:           totalIncome - totalExpense,
		IncomeByCategory:  incomeStats,
		ExpenseByCategory: expenseStats,
	}, nil
}

func (r *accounterFileRepo) GetPeriodStats(ctx context.Context, filter *biz.PeriodStatsFilter) (*biz.PeriodStats, error) {
	r.storage.mutex.RLock()
	defer r.storage.mutex.RUnlock()

	// 按时间段分组统计数据
	periodStats := make(map[string]*biz.PeriodData)
	var totalIncome, totalExpense float64

	for _, item := range r.storage.data {
		// Apply user filter
		if filter.UserID != 0 && item.UserID != filter.UserID {
			continue
		}

		// 根据时间段类型生成时间段名称
		var periodName string
		switch filter.PeriodType {
		case v1.PeriodType_MONTHLY:
			// 按月统计
			if filter.Year != 0 && item.Date.Year() != int(filter.Year) {
				continue
			}
			if filter.Month != 0 && item.Date.Month() != time.Month(filter.Month) {
				continue
			}
			periodName = fmt.Sprintf("%d年%d月", item.Date.Year(), item.Date.Month())
		case v1.PeriodType_YEARLY:
			// 按年统计
			if filter.Year != 0 && item.Date.Year() != int(filter.Year) {
				continue
			}
			periodName = fmt.Sprintf("%d年", item.Date.Year())
		case v1.PeriodType_WEEKLY:
			// 按周统计
			if filter.Year != 0 && item.Date.Year() != int(filter.Year) {
				continue
			}
			if filter.Week != 0 {
				// 计算指定年份的第N周
				yearStart := time.Date(int(filter.Year), 1, 1, 0, 0, 0, 0, time.UTC)
				weekStart := yearStart.AddDate(0, 0, (int(filter.Week)-1)*7)
				weekEnd := weekStart.AddDate(0, 0, 6)
				
				if item.Date.Before(weekStart) || item.Date.After(weekEnd) {
					continue
				}
			}
			year, week := item.Date.ISOWeek()
			periodName = fmt.Sprintf("%d年第%d周", year, week)
		default:
			continue
		}

		// 获取或创建时间段统计
		if stat, exists := periodStats[periodName]; exists {
			stat.TransactionCount++
		} else {
			periodStats[periodName] = &biz.PeriodData{
				PeriodName:       periodName,
				Income:           0,
				Expense:          0,
				Balance:          0,
				TransactionCount: 0,
			}
		}

		// 累计统计数据
		if item.Type == int32(v1.Type_Income) {
			periodStats[periodName].Income += item.Amount
			totalIncome += item.Amount
		} else if item.Type == int32(v1.Type_Expense) {
			periodStats[periodName].Expense += item.Amount
			totalExpense += item.Amount
		}
	}

	// 计算每个时间段的余额并转换为切片
	var periods []*biz.PeriodData
	for _, stat := range periodStats {
		stat.Balance = stat.Income - stat.Expense
		periods = append(periods, stat)
	}

	return &biz.PeriodStats{
		Periods:      periods,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		TotalBalance: totalIncome - totalExpense,
	}, nil
}
