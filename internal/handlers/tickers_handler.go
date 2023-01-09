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

type TickerInterface interface {
}

type TickerDatas struct {
	logger *zap.Logger
}

var _ TickerInterface = &TickerDatas{}

type Ticker struct {
	Pagination Pagination   `json:"pagination"`
	Ticker     []TickerData `json:"data"`
}

type TickerData struct {
	Name          string        `json:"name"`
	Symbol        string        `json:"symbol"`
	StockExchange StockExchange `json:"stock_exchange"`
}

type StockExchange struct {
	Name        string `json:"name"`
	Acronym     string `json:"acronym"`
	MIC         string `json:"mic"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	Website     string `json:"website"`
}

var _ TickerInterface = &Ticker{}
var ErrTickerInfo = "Error fetching stock tickers"

func NewTicker(logger *zap.Logger) *TickerDatas {
	return &TickerDatas{
		logger: logger,
	}
}

func (ticker *TickerDatas) GetStockTickers(ctx context.Context) (*Ticker, error) {
	if err := godotenv.Load(); err != nil {
		ticker.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrTickerInfo)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("STOCK_TICKER"), nil)
	if err != nil {
		ticker.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrTickerInfo)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		ticker.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrTickerInfo)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		ticker.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrTickerInfo)
	}
	var val Ticker
	if err := json.Unmarshal(data, &val); err != nil {
		ticker.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrTickerInfo)
	}
	return &val, nil
}
