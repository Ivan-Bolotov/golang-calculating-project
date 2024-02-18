package servers

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
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

type Server struct {
	id                 int
	Port               int `json:"port"`
	AmountOfGoroutines int `json:"amount"`
	state              string
	lastPing           time.Time
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

func (s *Server) Send(data []byte) bool {
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
			if err != nil || id >= len(expressions) {
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
					fmt.Fprintln(w, "not valid expression.")
					return
				}
				expressions = append(expressions, expression)
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
						server.Send(res)
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
		strings.Contains(expression.exp, ",") {
		return false
	} else if _, err = strconv.Atoi(expression.exp); err == nil {
		return false
	} else if _, err = strconv.ParseFloat(expression.exp, 64); err == nil {
		return false
	}
	return true
}