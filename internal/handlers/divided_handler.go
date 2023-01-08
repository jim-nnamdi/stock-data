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

type DividendInterface interface {
	GetDividends(ctx context.Context, symbol string) (*Dividend, error)
}

var _ DividendInterface = &Dividend{}
var ErrDividendError = "Error fetching dividend data"

type Dividend struct {
	logger     *zap.Logger
	Pagination Pagination     `json:"pagination"`
	Data       []DividendData `json:"data"`
}

type DividendData struct {
	Date     string  `json:"date"`
	Dividend float64 `json:"dividend"`
	Symbol   string  `json:"symbol"`
}

func NewDividend(logger *zap.Logger, pagination Pagination, data []DividendData) chan *Dividend {
	out := make(chan *Dividend)
	out <- &Dividend{
		logger:     logger,
		Pagination: pagination,
		Data:       data,
	}
	return out
}

func (dividend *Dividend) GetDividends(ctx context.Context, symbol string) (*Dividend, error) {
	if err := godotenv.Load(); err != nil {
		dividend.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New("unable to fetch env variables")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("STOCK_DIVIDENDS")+"&symbols="+symbol, nil)
	if err != nil {
		dividend.logger.Debug(err.Error(), zap.Error(err))
		return nil, errors.New(ErrDividendError)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		dividend.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrDividendError)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		dividend.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrDividendError)
	}
	var val Dividend
	if err := json.Unmarshal(data, &val); err != nil {
		dividend.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrDividendError)
	}
	return &val, nil
}
