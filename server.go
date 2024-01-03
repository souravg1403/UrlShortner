package main

//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
import (
	//"database/sql"

	"fmt"
	"log"
	"urlshortener/handlers"

	_ "github.com/lib/pq"

	//"log"

	"net/http"
	"os"
	"os/exec"
	"sync"
)

var port string = ":5000"

// func (h *Handler)

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

func main() {
	var wg sync.WaitGroup
	handler := handlers.NewHandler()
	wg.Add(1)
	go func() {

		defer wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/set", handler.SetURL)
		mux.HandleFunc("/", handler.FetchUrl)

		fmt.Println("Server is running on port " + port)

		server := &http.Server{
			Addr:    ":5000",
			Handler: mux,
		}

		log.Fatal(server.ListenAndServe())
	}()

	// wg.Add(1)
	// go servicesUp(&wg)

	wg.Wait()

	// postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)

	// fmt.Println(postgresqlDbInfo)
}
