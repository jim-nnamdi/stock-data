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

	historical    = handlers.NewHistoricalData(logger)
	ErrHistorical = "cannot fetch historical data"

	tickers    = handlers.NewTicker(logger)
	ErrTickers = "error fetching tickers!"
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

func HistoricalData(w http.ResponseWriter, r *http.Request) {
	symbol := r.FormValue("symbol")
	datefrom := r.FormValue("datefrom")
	dateto := r.FormValue("dateto")
	historicalData, err := historical.GetHistoricalEndOfDayData(r.Context(), symbol, datefrom, dateto)
	if err != nil {
		logger.Debug(err.Error(), zap.Any("error", err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(historicalData)
}

func GetStockTickers(w http.ResponseWriter, r *http.Request) {
	tickers, err := tickers.GetStockTickers(r.Context())
	if err != nil {
		logger.Debug(err.Error(), zap.Any("error", err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickers)
}
