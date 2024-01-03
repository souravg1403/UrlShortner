package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"urlshortener/store"

	"github.com/redis/go-redis/v9"
)

const (
	HOST          = "localhost"
	POSTGRES_PORT = 5432
	USERNAME      = "postgres"
	PASSWORD      = "sourav404"
	DATABASE      = "URL_Shortner"
)

type Handler struct {
	KeyValueStore store.KeyValueStore
}

func NewHandler() Handler {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", HOST),
		Password: "",
		DB:       0,
	})

	connStr := "user=" + USERNAME + " password=" + PASSWORD + " dbname=" + DATABASE + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Problem establishing postgres connection: ", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Problem checking postgres connection: ", err.Error())
	}

	store := store.NewKeyValueStore(db, redisClient)

	handler := &Handler{
		KeyValueStore: store,
	}

	return *handler
}
