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
