package data

import (
	"accounter/internal/conf"
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {

	db, err := NewGormDB(c, logger)
	if err != nil {
		return nil, nil, err
	}

	redisClient, cleanupRedis, err := NewRedisClient(c, logger)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		cleanupRedis()
	}

	return &Data{
		db:    db,
		redis: redisClient,
	}, cleanup, nil
}

func NewRedisClient(conf *conf.Data, logger log.Logger) (*redis.Client, func(), error) {
	client := redis.NewClient(&redis.Options{
		Addr: conf.Redis.Addr,
	})
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if err := client.Close(); err != nil {
			log.NewHelper(logger).Error("failed to close redis client: ", err)
		}
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.NewHelper(logger).Error("failed to ping redis: ", err)
		return nil, cleanup, err
	}

	return client, cleanup, nil
}

func NewGormDB(conf *conf.Data, logger log.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{})
	if err != nil {
		log.NewHelper(logger).Error("failed to open mysql: ", err)
		return nil, err
	}
	return db, nil
}
