package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"server/eval"
	"server/storage"
	"time"
)

type HttpApi struct {
	evaluator eval.ExpressionEvaluator
	storage   storage.HistoryStorage
}

func (a *HttpApi) Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/calculate", a.calculate).Methods(http.MethodPost)
	router.HandleFunc("/history", a.getHistory).Methods(http.MethodGet)

	return router
}

func respondWithJSON(w http.ResponseWriter, object interface{}, status int) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if object != nil { // non-empty body
		if err := json.NewEncoder(w).Encode(object); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	}
	return nil
}

type responseError struct {
	Error string
}

type responseResult struct {
	Result string
}

type requestCalculate struct {
	Expression string
}

func (a *HttpApi) calculate(w http.ResponseWriter, r *http.Request) {
	var calc requestCalculate
	if err := json.NewDecoder(r.Body).Decode(&calc); err != nil {
		if err := respondWithJSON(w, responseError{
			fmt.Sprintf("Invalid request: %v", err.Error())},
			http.StatusBadRequest); err != nil {
			return
		}
		return
	}

	result, err := a.evaluator.Evaluate(calc.Expression)
	if err != nil {
		if err := respondWithJSON(w, responseError{err.Error()}, http.StatusOK); err != nil {
			return
		}
		return
	}
	if err := a.storage.StoreCalculation(storage.Calculation{
		Expression: calc.Expression,
		Result:     result,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := respondWithJSON(w, responseResult{result}, http.StatusOK); err != nil {
		return
	}
}

func (a *HttpApi) getHistory(w http.ResponseWriter, r *http.Request) {
	result, err := a.storage.GetHistory()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := respondWithJSON(w, result, http.StatusOK); err != nil {
		return
	}
}

func main() {
	connStr := "user=postgres password=123 host=db dbname=postgres sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	a := &HttpApi{evaluator: &eval.SmartEvaluator{}, storage: storage.New(conn)}
	server := http.Server {
		Addr: ":8081",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler: a.Router(),
	}
    println("Running server")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
