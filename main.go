package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jim-nnamdi/stock-data/internal/handlers"
	"go.uber.org/zap"
)

var (
	logger      = zap.NewNop()
	endofday    = handlers.NewEOD(logger)
	ErrEndOfDay = "cannot fetch end of day from client"
)

func EODentry(w http.ResponseWriter, r *http.Request) {
	symbol := r.FormValue("symbol")
	eodresult, err := endofday.GetEndOfDayData(r.Context(), symbol)
	if err != nil {
		logger.Debug(ErrEndOfDay, zap.Any(ErrEndOfDay, err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eodresult)
}

func main() {
	msg := make(chan string)
	go func() {
		msg <- "Hello nyse & nasdaq"
	}()
	mx := <-msg
	fmt.Printf("hello stock market!: %v", mx)
}
