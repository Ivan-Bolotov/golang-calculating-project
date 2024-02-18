package servers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func StartNewHttpComputingServer(port int) {
	envVars, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	amount, err := strconv.Atoi(envVars["AMOUNT_OF_GOROUTINES"])
	if err != nil {
		log.Fatal("Invalid amount of goroutines in .env file")
	}

	rootServerPort, err := strconv.Atoi(envVars["ROOT_SERVER_PORT"])
	if err != nil {
		log.Fatal("Port in .env file must be integer")
	}

	fmt.Printf("%d goroutines required.\n", amount)

	muxHandler := mux.NewRouter()
	setRoutesForComputing(muxHandler)

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: muxHandler}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/add_server", rootServerPort), strings.NewReader(fmt.Sprintf(`{"port": %d, "amount": %d}`, port, amount)))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	if err = server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server")
	}
}

func setRoutesForComputing(handler *mux.Router) {
	handler.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			return
		}
	})
}
