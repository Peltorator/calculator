package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

func (a *HttpApi) calculate(w http.ResponseWriter, r *http.Request) {
	var expr string
	if err := json.NewDecoder(r.Body).Decode(&expr); err != nil {
		if err := respondWithJSON(w, responseError{
			fmt.Sprintf("Request is not a string: %v", err.Error())},
			http.StatusBadRequest); err != nil {
			return
		}
		return
	}

	result, err := a.evaluator.Evaluate(expr)
	if err != nil {
		if err := respondWithJSON(w, responseError{err.Error()}, http.StatusOK); err != nil {
			return
		}
		return
	}
	if err := a.storage.StoreCalculation(storage.Calculation{
		Expression: expr,
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
	a := &HttpApi{evaluator: &eval.IncrementalFakeEvaluator{}, storage: &storage.InMemoryHistoryStorage{
		Calculations: make([]storage.Calculation, 0),
	}}
	server := http.Server {
		Addr: ":8080",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler: a.Router(),
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}