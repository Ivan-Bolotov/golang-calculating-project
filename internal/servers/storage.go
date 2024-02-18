package servers

import (
	"encoding/json"
	"fmt"
	"github.com/Ivan-Bolotov/golang-calculating-project/pkg/functions"
	"github.com/Knetic/govaluate"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Expression struct {
	id     int
	exp    string
	result interface{}
	state  string
}

type StringExpression struct {
	String string `json:"expression"`
}

type ExpressionResult struct {
	Id     int         `json:"id"`
	Result interface{} `json:"result"`
}

type Server struct {
	id                  int
	Port                int `json:"port"`
	AmountOfGoroutines  int `json:"amount"`
	state               string
	lastPing            time.Time
	countingExpressions []int
}

func (s *Server) Ping() {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/ping", s.Port), strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	s.lastPing = time.Now()
	if err != nil || resp.StatusCode != http.StatusOK {
		s.state = "inactive"
		return
	}
	s.state = "active"
}

func (s *Server) SendOperationTime(data []byte) bool {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/set_operation_time", s.Port), strings.NewReader(string(data)))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	s.lastPing = time.Now()
	if err != nil || resp.StatusCode != http.StatusOK {
		s.state = "inactive"
		return false
	}
	return true
}

func (s *Server) SendExpression(data []byte) bool {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/expression", s.Port), strings.NewReader(string(data)))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	s.lastPing = time.Now()
	if err != nil || resp.StatusCode != http.StatusOK {
		s.state = "inactive"
		return false
	}
	return true
}

var (
	expressions        []Expression
	computingResources []Server
)

func StartNewHttpStorageServer() {
	envVars, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rootServerPort, err := strconv.Atoi(envVars["ROOT_SERVER_PORT"])
	if err != nil {
		log.Fatal("Port in .env file must be integer")
	}

	muxHandler := mux.NewRouter()
	setRoutesForStorage(muxHandler)
	server := &http.Server{Addr: fmt.Sprintf(":%d", rootServerPort), Handler: muxHandler}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server")
	}
}

func setRoutesForStorage(handler *mux.Router) {
	handler.HandleFunc("/expression/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			vars := mux.Vars(r)
			id, err := strconv.Atoi(vars["id"])
			if err != nil || id >= len(expressions) || id < 0 {
				fmt.Fprintln(w, "not valid id.")
				return
			}
			expression := expressions[id]
			var dict = map[string]interface{}{
				"state": expression.state,
				"res":   expression.result,
				"exp":   expression.exp,
			}
			data, err := json.Marshal(dict)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(w, string(data))
			// TODO: лезем в БД или в Redis и отправляем клиенту оттуда state и res по ID
		}
	})
	handler.HandleFunc("/{path}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		switch vars["path"] {
		case "expression":
			if r.Method == http.MethodPost {
				var exp StringExpression
				res, _ := io.ReadAll(r.Body)
				err := json.Unmarshal(res, &exp)
				if err != nil {
					panic(err)
				}
				var expression = Expression{
					id:     len(expressions),
					exp:    strings.Join(strings.Split(exp.String, " "), ""),
					state:  "pending",
					result: nil,
				}
				if !isValidExpression(expression) {
					fmt.Fprintln(w, "invalid expression.")
					return
				}
				if !isUniqueExpression(expressions, expression) {
					fmt.Fprintln(w, "not unique expression.")
					return
				}
				expressions = append(expressions, expression)

				var lensOfCountingExpressionsSlices []int
				for _, server := range computingResources {
					server.Ping()
					if server.state == "active" {
						lensOfCountingExpressionsSlices = append(lensOfCountingExpressionsSlices, len(server.countingExpressions))
					}
				}
				m, _ := functions.MinMax(lensOfCountingExpressionsSlices)
				for _, server := range computingResources {
					if server.state == "active" && len(server.countingExpressions) == m {
						server.SendExpression([]byte(fmt.Sprintf(`{"id": %d, "expression": "%s"}`, expression.id, expression.exp)))
					}
				}
				var dict = map[string]interface{}{
					"id":    expression.id,
					"state": expression.state,
					"res":   expression.result,
				}
				data, err := json.Marshal(dict)
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w, string(data))
				// TODO: отдаём выражение на вычисление свободному серверу и кладём его в БД

			}
		case "expressions":
			if r.Method == http.MethodGet {
				var arr []map[string]interface{}
				for _, expression := range expressions {
					var dict = map[string]interface{}{
						"id":    expression.id,
						"state": expression.state,
						"res":   expression.result,
						"exp":   expression.exp,
					}
					arr = append(arr, dict)
				}
				data, err := json.Marshal(arr)
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w, string(data))
				// TODO: идём в базу и достаём все записи
			}
		case "computing_resources":
			if r.Method == http.MethodGet {
				var arr []map[string]interface{}
				for _, server := range computingResources {
					server.Ping()
					var dict = map[string]interface{}{
						"id":       server.id,
						"state":    server.state,
						"lastping": server.lastPing.String(),
					}
					arr = append(arr, dict)
				}
				data, err := json.Marshal(arr)
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w, string(data))
				// TODO: проверяем подключение ко всем серверам (горутинам вычислителя) и возвращаем результат
			}
		case "set_operation_time":
			if r.Method == http.MethodPost {
				res, _ := io.ReadAll(r.Body)
				for _, server := range computingResources {
					server.Ping()
					if server.state == "active" {
						server.SendOperationTime(res)
					}
				}
				fmt.Fprintln(w, "ok.")
				// TODO: Доделать!!! [ DONE ]
				// TODO: отправляем на сервера-вычислители время выполнения данного вычисления
			}
		case "add_server":
			if r.Method == http.MethodPost {
				var server Server
				res, _ := io.ReadAll(r.Body)
				err := json.Unmarshal(res, &server)
				if err != nil {
					panic(err)
				}

				server.id = len(computingResources)
				server.state = "active"
				server.lastPing = time.Now()

				computingResources = append(computingResources, server)

				var id = map[string]int{"id": server.id}
				data, err := json.Marshal(id)
				if err != nil {
					panic(err)
				}
				fmt.Fprintln(w, string(data))
				// TODO: добавляем новый сервак
			}
		case "result":
			if r.Method == http.MethodPost {
				var expResult ExpressionResult
				res, _ := io.ReadAll(r.Body)
				err := json.Unmarshal(res, &expResult)
				if err != nil {
					panic(err)
				}
				expressions[expResult.Id].result = expResult.Result
				expressions[expResult.Id].state = "resolved"
				for _, server := range computingResources {
					if slices.Contains(server.countingExpressions, expResult.Id) {
						for i, val := range server.countingExpressions {
							if val == expResult.Id {
								server.countingExpressions = append(server.countingExpressions[:i], server.countingExpressions[i+1:]...)
							}
						}
					}
				}
				fmt.Fprintln(w, "ok.")
				// TODO: изменение данных в БД
			}
		}
	})
}

func isValidExpression(expression Expression) bool {
	_, err := govaluate.NewEvaluableExpression(expression.exp)
	if err != nil {
		return false
	} else if strings.Contains(expression.exp, "**") ||
		strings.Contains(expression.exp, "%") ||
		strings.Contains(expression.exp, "|") ||
		strings.Contains(expression.exp, ",") ||
		strings.Contains(expression.exp, "/0") {
		return false
	} else if _, err = strconv.Atoi(expression.exp); err == nil {
		return false
	} else if _, err = strconv.ParseFloat(expression.exp, 64); err == nil {
		return false
	}
	return true
}

func isUniqueExpression(expressions []Expression, expression Expression) bool {
	for _, exp := range expressions {
		if exp.exp == expression.exp {
			return false
		}
	}
	return true
}
