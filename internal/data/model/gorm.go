package model

import (
	"time"
)

// AccounterCategory 交易分类表，用于记录各类支出或收入的分类
type AccounterCategory struct {
	CategoryID   int       `gorm:"column:category_id;primaryKey;autoIncrement" json:"category_id"`                                       // 分类主键ID，自增
	CategoryName string    `gorm:"column:category_name;type:varchar(50);not null" json:"category_name"`                                  // 分类名称，如餐饮、交通、工资等
	ParentID     *int      `gorm:"column:parent_id;type:int" json:"parent_id"`                                                           // 父级分类ID，用于多级分类，自关联到本表category_id，可为空
	Type         *int8     `gorm:"column:type;type:tinyint" json:"type"`                                                                 // 分类类型：0-支出，1-收入；可选字段，若不区分可不使用
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"created_at"`                // 分类创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;not null;autoUpdateTime" json:"updated_at"` // 分类更新时间
}

// TableName 设置表名
func (AccounterCategory) TableName() string {
	return "accounter_categories"
}

// AccounterTransaction 交易明细表，记录每笔收入或支出信息
type AccounterTransaction struct {
	TransactionID   int64     `gorm:"column:transaction_id;primaryKey;autoIncrement" json:"transaction_id"`                                 // 交易主键ID，自增
	UserID          int64     `gorm:"column:user_id;type:bigint;not null" json:"user_id"`                                                   // 用户ID, 关联users.user_id
	CategoryID      int       `gorm:"column:category_id;type:int;not null" json:"category_id"`                                              // 交易所属分类ID，外键关联categories.category_id
	CurrencyID      int       `gorm:"column:currency_id;type:int;not null" json:"currency_id"`                                              // 使用的币种ID，外键关联currencies.currency_id
	TransactionType int8      `gorm:"column:transaction_type;type:tinyint;not null" json:"transaction_type"`                                // 交易类型：0-支出，1-收入
	Amount          float64   `gorm:"column:amount;type:decimal(18,5);not null" json:"amount"`                                              // 交易金额，一般保留两位小数
	TransactionDate time.Time `gorm:"column:transaction_date;type:datetime;not null" json:"transaction_date"`                               // 交易实际发生时间
	Note            *string   `gorm:"column:note;type:varchar(255)" json:"note"`                                                            // 交易备注信息，如“早餐”、“地铁费”等
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"created_at"`                // 记录创建时间
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;not null;autoUpdateTime" json:"updated_at"` // 记录更新时间
}

// TableName 设置表名
func (AccounterTransaction) TableName() string {
	return "accounter_transactions"
}

// Currency 币种信息表，用于维护可用的货币类型
type Currency struct {
	CurrencyID     int       `gorm:"column:currency_id;primaryKey;autoIncrement" json:"currency_id"`                                       // 币种主键ID，自增
	CurrencyCode   string    `gorm:"column:currency_code;type:varchar(10);not null;unique" json:"currency_code"`                           // 币种代码，例如 CNY, USD, EUR
	CurrencyName   string    `gorm:"column:currency_name;type:varchar(50);not null" json:"currency_name"`                                  // 币种名称，例如 人民币, 美元, 欧元
	CurrencySymbol *string   `gorm:"column:currency_symbol;type:varchar(10)" json:"currency_symbol"`                                       // 币种符号，例如 ¥, $, €
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"created_at"`                // 记录创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;not null;autoUpdateTime" json:"updated_at"` // 记录更新时间
}

// TableName 设置表名
func (Currency) TableName() string {
	return "currencies"
}

// User 用户信息表
type User struct {
	UserID       int       `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`                                                 // 主键ID，自增
	Username     string    `gorm:"column:username;type:varchar(50);not null;unique" json:"username"`                                       // 用户名
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null" json:"password_hash"`                                   // 密码hash值
	CreateTime   time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"create_time"`                // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;not null;autoUpdateTime" json:"update_time"` // 更新时间
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}
