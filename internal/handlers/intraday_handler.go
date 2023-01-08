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

type IntraDayInterface interface {
	GetIntraDay(ctx context.Context, symbol string) (*IntraDay, error)
	LatestIntraDay(ctx context.Context, symbol string) (*IntraDay, error)
	RealTimeIntraDayUpdate(ctx context.Context, symbol string, interval string) (*IntraDay, error)
}

var _ IntraDayInterface = &IntraDay{}

var (
	ErrIntraDayOutputStream = "cannot feed stock data into output stream"
)

type IntraDay struct {
	logger     *zap.Logger
	Pagination struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Count  int `json:"count"`
		Total  int `json:"total"`
	} `json:"pagination"`
	Data []struct {
		Open     float64 `json:"open"`
		High     float64 `json:"high"`
		Low      float64 `json:"low"`
		Last     float64 `json:"last"`
		Close    float64 `json:"close"`
		Volume   float64 `json:"volume"`
		Date     string  `json:"date"`
		Symbol   string  `json:"symbol"`
		Exchange string  `json:"exchange"`
	} `json:"data"`
}

func (intraday *IntraDay) GetIntraDay(ctx context.Context, symbol string) (*IntraDay, error) {
	if err := godotenv.Load(); err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("STOCK_BASE_INTRADAY")+"&symbols="+symbol, nil)
	if err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()
	dataBytes, err := io.ReadAll(res.Body)
	if err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	var val IntraDay
	if err := json.Unmarshal(dataBytes, &val); err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrIntraDayOutputStream)
	}
	return &val, nil
}

func (intraday *IntraDay) LatestIntraDay(ctx context.Context, symbol string) (*IntraDay, error) {
	if err := godotenv.Load(); err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("STOCK_BASE_INTRADAY_LATEST")+"&symbols="+symbol, nil)
	if err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()
	dataBytes, err := io.ReadAll(res.Body)
	if err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	var val IntraDay
	if err := json.Unmarshal(dataBytes, &val); err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrIntraDayOutputStream)
	}
	return &val, nil
}

// the interval parameter will be analysed
// first taken in as time.Duration and then
// parsed into a time.Time type <database-use>
func (intraday *IntraDay) RealTimeIntraDayUpdate(ctx context.Context, symbol string, interval string) (*IntraDay, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("STOCK_BASE_EOD")+"&symbols="+symbol+"&interval="+interval, nil)
	if err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrIntraDayOutputStream)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		intraday.logger.Debug(err.Error(), zap.Error(err))
		return nil, errors.New(ErrIntraDayOutputStream)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrIntraDayOutputStream)
	}
	var val IntraDay
	if err := json.Unmarshal(data, &val); err != nil {
		intraday.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrIntraDayOutputStream)
	}
	return &val, nil
}
