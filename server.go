package main

//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
import (
	"./UrlShortner/helpers"
	"encoding/json"
	"fmt"
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

func (redis *RedisConnection) urlShortening(u URL) string {
	randomKey := getRandomString()
	check := redis.isExist(randomKey)
	for check {
		randomKey = getRandomString()
		check = redis.isExist(randomKey)
	}

	result, err := redis.setKey(randomKey, u.WebAddress, 0)
	if err != nil {
		fmt.Println("Error:", err)
		return "error"
	} else {
		fmt.Println(result)
		return randomKey
	}
}

func (redis *RedisConnection) GetUrl(w http.ResponseWriter, r *http.Request) {
	var data URL
	decoder := json.NewDecoder(r.Body)
	myData := decoder.Decode(&data)
	fmt.Println(myData)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println(err)
	}

	value := redis.urlShortening(data)
	shortenedUrl := "http://localhost" + port + "/" + value
	shorturl := shortURL{
		ShortURL: shortenedUrl,
	}

	dataBody, _ := json.Marshal(shorturl)
	w.Write(dataBody)
}

func (redis *RedisConnection) FetchUrl(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) >= 2 {
		extractedPart := parts[1]
		fmt.Println(extractedPart)
		fmt.Printf("type: %T", extractedPart)

		value, err := redis.getValue(extractedPart)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, value, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	var wg sync.WaitGroup
	handler := RedisConnectionHandler()
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
