package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jim-nnamdi/stock-data/internal/entrypoint"
)

func EndOfDayInfo(w http.ResponseWriter, r *http.Request) {
	entrypoint.EndOfDayInfo(w, r)
}
func LatestEndOfDayInfo(w http.ResponseWriter, r *http.Request) {
	entrypoint.EndOfDayInfo(w, r)
}
func Dividends(w http.ResponseWriter, r *http.Request) {
	entrypoint.Dividends(w, r)
}

func main() {
	route := mux.NewRouter()
	route.HandleFunc("/eod", EndOfDayInfo)
	route.HandleFunc("/eod/latest", LatestEndOfDayInfo)
	route.HandleFunc("/dividends", Dividends)

	msg := make(chan string)
	go func() {
		msg <- "Hello nyse & nasdaq"
	}()
	mx := <-msg
	fmt.Printf("hello stock market!: %v", mx)
}
