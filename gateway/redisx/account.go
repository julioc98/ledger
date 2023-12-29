package redisx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/julioc98/ledger/domain/entities"
	"github.com/julioc98/ledger/gateway/pg"
	"github.com/redis/go-redis/v9"
)

const (
	cachePrefix = "ledger:"
	cacheTTL    = 5 * time.Second // Cache expiration time
)

// AccountRedisRepositoryDecorator is a repository that uses Redis as a database.
type AccountRedisRepositoryDecorator struct {
	redisClient *redis.Client
	dbRepo      *pg.AccountPgxRepository
}

// NewAccountRedisRepositoryDecorator returns a new AccountRedisRepositoryDecorator.
func NewAccountRedisRepositoryDecorator(redisClient *redis.Client, dbRepo *pg.AccountPgxRepository) *AccountRedisRepositoryDecorator {
	return &AccountRedisRepositoryDecorator{redisClient, dbRepo}
}

func (r *AccountRedisRepositoryDecorator) TransfersHistory(ctx context.Context, account string) ([]entities.Entry, error) {
	// Cache key for the transfer history
	cacheKey := cachePrefix + "transfers:" + account

	// Try to get the transfer history from the Redis cache
	cacheValue, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		var history []entities.Entry
		if err := json.Unmarshal([]byte(cacheValue), &history); err == nil {
			return history, nil
		}
	}

	// Cache miss, get the transfer history from the database
	history, err := r.dbRepo.TransfersHistory(ctx, account)
	if err != nil {
		return nil, err
	}

	// Save the transfer history in the cache as a JSON string
	cacheValueBytes, err := json.Marshal(history)
	if err != nil {
		return nil, err
	}

	cacheValue = string(cacheValueBytes) // Convert bytes to string

	err = r.redisClient.Set(ctx, cacheKey, cacheValue, cacheTTL).Err()
	if err != nil {
		return nil, err
	}

	return history, nil
}
