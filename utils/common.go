package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func GetRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	randomKey := make([]byte, 5)
	for i := range randomKey {
		randomKey[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(randomKey)
}

func JSONError(w http.ResponseWriter, err interface{}, code int) http.ResponseWriter {
	fmt.Printf("Error Occurred: %v\n", err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
	return w
}
