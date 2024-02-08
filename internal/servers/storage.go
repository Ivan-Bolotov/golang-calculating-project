package servers

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func StartNewHttpStorageServer() {
	muxHandler := mux.NewRouter()
	setRoutes(muxHandler)
	server := &http.Server{Addr: ":8080", Handler: muxHandler}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("failed to start server.")
	}
}

func setRoutes(handler *mux.Router) {
	handler.HandleFunc("/expression/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// vars := mux.Vars(r)
			// id := vars["id"]
			// TODO: лезем в БД или в Redis и отправляем клиенту оттуда state и res по ID
		}
	})
	handler.HandleFunc("/{path}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		switch vars["path"] {
		case "expression":
			if r.Method == http.MethodPost {
				res, _ := io.ReadAll(r.Body)
				body := string(res)
				fmt.Println(body)
				// TODO: отдаём выражение на вычисление свободному серверу и кладём его в БД

			}
		case "expressions":
			if r.Method == http.MethodGet {
				// TODO: идём в базу и достаём все записи
			}
		case "computing_resources":
			if r.Method == http.MethodGet {
				// TODO: проверяем подключение ко всем серверам (горутинам вичислителя) и возвращаем результат
			}
		case "set_operation_time":
			if r.Method == http.MethodPost {
				res, _ := io.ReadAll(r.Body)
				body := string(res)
				fmt.Println(body)
				// TODO: отправляем на сервера-вычислители время выполнения данного вичисления
			}
		}
	})
}
