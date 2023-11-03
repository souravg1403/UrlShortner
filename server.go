package main

//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
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type URL struct {
	WebAddress string `json:"url"`
}

type shortURL struct {
	ShortURL string `json:"shorturl"`
}

var port string = ":5000"

func servicesUp(wg *sync.WaitGroup) {

	defer wg.Done()
	fmt.Println("Services getting up")
	cmd := exec.Command("docker", "compose", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		fmt.Println("Services getting up")
	}
}

func getRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	randomKey := make([]byte, 5)
	for i := range randomKey {
		randomKey[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(randomKey)
}

func (h *Handler) urlShortening(u URL) string {

	//Connection creation with redis
	ctx := context.TODO()
	// fmt.Println(u)
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       0,
	// })

	randomKey := getRandomString()
	check, _ := h.RedisClient.Exists(ctx, randomKey).Result()
	for check == 1 {
		randomKey = getRandomString()
		check, _ = h.RedisClient.Exists(ctx, randomKey).Result()
	}

	fmt.Println(randomKey)

	result := h.RedisClient.Set(ctx, randomKey, u.WebAddress, 0)
	err := result.Err()
	if err != nil {
		fmt.Println("Error:", err)
		return "error"
	}
	return randomKey
}

// Handler function
func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) {

	var data URL
	decoder := json.NewDecoder(r.Body)
	myData := decoder.Decode(&data)
	fmt.Println(myData)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println(err)
	}

	value := h.urlShortening(data)
	shortenedUrl := "http://localhost" + port + "/" + value
	shorturl := shortURL{
		ShortURL: shortenedUrl,
	}

	dataBody, _ := json.Marshal(shorturl)
	w.Write(dataBody)

}

func (h *Handler) FetchUrl(w http.ResponseWriter, r *http.Request) {

	//Connection creation with redis
	ctx := context.TODO()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) >= 2 {
		extractedPart := parts[1]
		fmt.Println(extractedPart)
		fmt.Printf("type: %T", extractedPart)

		searchKey := h.RedisClient.Get(ctx, extractedPart)
		http.Redirect(w, r, searchKey.Val(), http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

type Handler struct {
	RedisClient *redis.Client
}

func NewHandler() Handler {
	fmt.Println("########")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	handler := Handler{
		RedisClient: client,
	}

	fmt.Printf("type: %T", client)
	return handler
}

func main() {
	var wg sync.WaitGroup
	handler := NewHandler()
	wg.Add(1)
	go func() {

		defer wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/get_url", handler.GetUrl)
		mux.HandleFunc("/", handler.FetchUrl)

		fmt.Println("Server is running on port " + port)

		server := &http.Server{
			Addr:    ":5000",
			Handler: mux,
		}

		log.Fatal(server.ListenAndServe())
	}()

	wg.Add(1)
	go servicesUp(&wg)

	wg.Wait()

}
