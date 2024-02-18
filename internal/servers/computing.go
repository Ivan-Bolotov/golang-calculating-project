package servers

import (
	"encoding/json"
	"fmt"
	"github.com/Ivan-Bolotov/golang-calculating-project/pkg/workerpool"
	"github.com/Knetic/govaluate"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	rootServerPort int
	workerPool     *workerpool.Pool
	operandsList   = []string{"+", "-", "/", "*"}
	operandsMap    = map[string]time.Duration{"+": 0, "-": 0, "/": 0, "*": 0}
)

func StartNewHttpComputingServer(port int) {
	// читаем файл с переменными окружения
	envVars, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// кол-во горутин
	amount, err := strconv.Atoi(envVars["AMOUNT_OF_GOROUTINES"])
	if err != nil {
		log.Fatal("Invalid amount of goroutines in .env file")
	}

	// порт главного сервера (оркестратора)
	rootServerPort, err = strconv.Atoi(envVars["ROOT_SERVER_PORT"])
	if err != nil {
		log.Fatal("Port in .env file must be integer")
	}

	fmt.Printf("%d goroutines required.\n", amount)

	// создаём WorkerPool
	workerPool = workerpool.New(amount, worker)

	// настраиваем обработку запросов
	muxHandler := mux.NewRouter()
	setRoutesForComputing(muxHandler)

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: muxHandler}
	client := &http.Client{}

	// даём о себе знать главному серверу
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
			fmt.Fprintln(w, "ok.")
		}
	})
	handler.HandleFunc("/set_operation_time", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var data = map[string]int{}
			res, _ := io.ReadAll(r.Body)
			err := json.Unmarshal(res, &data)
			if err != nil {
				log.Fatal(err)
			}
			for key, value := range data {
				if slices.Contains(operandsList, key) {
					operandsMap[key] = time.Duration(value * int(math.Pow(10, 9)))
				}
			}
			fmt.Fprintln(w, "ok.")
		}
	})
	handler.HandleFunc("/expression", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var data = map[string]interface{}{}
			res, _ := io.ReadAll(r.Body)
			err := json.Unmarshal(res, &data)
			if err != nil {
				log.Fatal(err)
			}
			workerPool.Add(data)
			fmt.Fprintln(w, "ok.")
		}
	})
}

func worker(exp map[string]interface{}) {
	var (
		expression = fmt.Sprintf("%v", exp["expression"])
		dur        time.Duration
	)
	id, err := strconv.Atoi(fmt.Sprintf("%v", exp["id"]))
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range strings.Split(expression, "") {
		if subDur, ok := operandsMap[val]; ok {
			dur += subDur
		}
	}
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		log.Fatal(err)
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(dur)
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/result", rootServerPort), strings.NewReader(fmt.Sprintf(`{"id": %d, "result": %f}`, id, result)))
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
}
