package entrypoint

import (
	"encoding/json"
	"net/http"

	"github.com/jim-nnamdi/stock-data/internal/handlers"
	"go.uber.org/zap"
)

var (
	logger      = zap.NewNop()
	endofday    = handlers.NewEODData(logger)
	ErrEndOfDay = "cannot fetch end of day from client"

	dividend    = handlers.NewDividendData(logger)
	ErrDividend = "cannot fetch dividends from client"
)

func EndOfDayInfo(w http.ResponseWriter, r *http.Request) {
	symbol := r.FormValue("symbol")
	eodresult, err := endofday.GetEndOfDayData(r.Context(), symbol)
	if err != nil {
		logger.Debug(ErrEndOfDay, zap.Any(ErrEndOfDay, err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eodresult)
}

func LatestEndOfDayInfo(w http.ResponseWriter, r *http.Request) {
	symbol := r.FormValue("symbol")
	eodresult, err := endofday.GetEndOfDayData(r.Context(), symbol)
	if err != nil {
		logger.Debug(ErrEndOfDay, zap.Any(ErrEndOfDay, err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eodresult)
}

func Dividends(w http.ResponseWriter, r *http.Request) {
	symbol := r.FormValue("symbol")
	dividends, err := dividend.GetDividends(r.Context(), symbol)
	if err != nil {
		logger.Debug(ErrDividend, zap.Any(ErrDividend, err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dividends)
}
