package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type HistoricalData interface {
	GetHistoricalEndOfDayData(ctx context.Context, symbol string, datefrom string, dateto string) (*Historical, error)
}

var _ HistoricalData = &Historical{}

var (
	ErrHistoricalOutputStream = "cannot feed historical data into output stream"
)

type Historical struct {
	logger     *zap.Logger
	Pagination struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Count  int `json:"count"`
		Total  int `json:"total"`
	} `json:"pagination"`
	Data []struct {
		Open      float64 `json:"open"`
		High      float64 `json:"high"`
		Low       float64 `json:"low"`
		Close     float64 `json:"close"`
		Volume    float64 `json:"volume"`
		AdjHigh   float64 `json:"adj_high"`
		AdjLow    float64 `json:"adj_low"`
		AdjClose  float64 `json:"adj_close"`
		AdjOpen   float64 `json:"adj_open"`
		AdjVolume float64 `json:"adj_volume"`
	} `json:"data"`
}

// Get Historical end of day data for a particular symbol
// the dates are taken in as string as time.duration and
// parsed as time.Time if used in a database.
func (historical *Historical) GetHistoricalEndOfDayData(ctx context.Context, symbol string, datefrom string, dateto string) (*Historical, error) {
	if err := godotenv.Load(); err != nil {
		historical.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("STOCK_BASE_EOD")+"&symbols="+symbol+"&date_from="+datefrom+"&date_to ="+dateto, nil)
	if err != nil {
		historical.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		historical.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()
	dataBytes, err := io.ReadAll(res.Body)
	if err != nil {
		historical.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	var val Historical
	if err := json.Unmarshal(dataBytes, &val); err != nil {
		historical.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrEODOutputStream)
	}
	return &val, nil
}
