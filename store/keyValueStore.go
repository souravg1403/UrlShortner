package store

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	DEFAULT_REDIS_EXPIRY_IN_SECONDS = 300
)

type KeyValueStore struct {
	Redis *redis.Client
	Pg    *sql.DB
}

func NewKeyValueStore(db *sql.DB, redisClient *redis.Client) KeyValueStore {
	store := &KeyValueStore{
		Redis: redisClient,
		Pg:    db,
	}

	return *store
}

func (store *KeyValueStore) SetKey(ctx context.Context, key, value string) (*string, error) {
	redisTx := store.Redis.TxPipeline()

	// Add Redis write operation
	redisTx.Set(ctx, key, value, DEFAULT_REDIS_EXPIRY_IN_SECONDS*time.Second)

	// Start a PostgreSQL transaction
	pgTx, err := store.Pg.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Add PostgreSQL write operation
	_, err = pgTx.Exec("INSERT INTO your_table (key, value) VALUES ($1, $2)", key, value)
	if err != nil {
		pgTx.Rollback()
		return nil, err
	}

	// Commit Postgres
	err = pgTx.Commit()
	if err != nil {
		return nil, err
	}

	// Commit both transactions
	_, _ = redisTx.Exec(ctx)

	return nil, nil
}

func (store *KeyValueStore) GetKey(ctx context.Context, key string) (*string, bool, error) {
	// Checking for key in redis
	keyExist, err := store.Redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, false, err
	}

	var responseVal string

	// If key exist in redis return the value and bool confirmation
	if keyExist == 1 {
		response := store.Redis.Get(ctx, key)
		responseVal = response.Val()
	} else {
		// If key doesn't exist in redis check in postgres. If in pg load into redis else return not found

		err := store.Pg.QueryRow("SELECT url FROM key_url WHERE key = $1", key).Scan(&responseVal)
		if err != nil {
			return nil, false, err
		}

		if responseVal != "" {
			response := store.Redis.Set(ctx, key, responseVal, time.Duration(DEFAULT_REDIS_EXPIRY_IN_SECONDS*time.Second))
			_, err := response.Result()
			if err != nil {
				return nil, false, err
			}
		}
	}

	if responseVal == "" {
		return nil, false, nil
	}

	return &responseVal, true, nil
}

// func (store *KeyValueStore) DeleteKeyFromRedis(ctx context.Context, key, value string) (*string, error) {
// 	keyExist, err := store.Redis.Exists(ctx, key).Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// If key exist in redis return the value and bool confirmation
// 	if keyExist == 1 {
// 		return nil, nil
// 	}

// 	return nil, nil
// }
