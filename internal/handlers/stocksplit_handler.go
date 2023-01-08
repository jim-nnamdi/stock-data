package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type SplitFactorInterface interface {
	GetSplitFactor(ctx context.Context, symbol string) (*SplitFactor, error)
}

var _ SplitFactorInterface = &SplitFactor{}

var (
	ErrSplitFactorOutputStream = "cannot feed historical data into output stream"
)

type SplitFactor struct {
	logger     *zap.Logger
	Pagination Pagination  `json:"pagination"`
	Data       []SplitData `json:"data"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
	Total  int `json:"total"`
}

type SplitData struct {
	SplitFactor float64 `json:"split_factor"`
	Symbol      string  `json:"symbol"`
	Date        string  `json:"date"`
}

func NewSplitFactor(logger *zap.Logger, pagination Pagination, data []SplitData) chan *SplitFactor {
	out := make(chan *SplitFactor)
	out <- &SplitFactor{
		logger:     logger,
		Pagination: pagination,
		Data:       data,
	}
	return out
}

func (sf *SplitFactor) GetSplitFactor(ctx context.Context, symbol string) (*SplitFactor, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("STOCK_SPLITFACTOR")+"&symbol="+symbol, nil)
	if err != nil {
		sf.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrSplitFactorOutputStream)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		sf.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrSplitFactorOutputStream)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		sf.logger.Debug(err.Error(), zap.Error(err))
		return nil, errors.New(ErrSplitFactorOutputStream)
	}
	var val SplitFactor
	if err := json.Unmarshal(data, &val); err != nil {
		sf.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrSplitFactorOutputStream)
	}
	return &val, nil
}
