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

type EndOfDay interface {
	GetEndOfDayData(ctx context.Context, symbol string) (*EOD, error)
	LatestEndOfDayData(ctx context.Context, symbol string) (*EOD, error)
}

var _ EndOfDay = &NewEOD{}

var (
	client             = &http.Client{}
	ErrEODOutputStream = "cannot feed stock data into output stream"
)

type NewEOD struct {
	logger *zap.Logger
}

type EOD struct {
	Pagination Pagination `json:"pagination"`
	Data       []EODData  `json:"data"`
}

type EODData struct {
	Open        float64 `json:"open"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Close       float64 `json:"close"`
	Volume      float64 `json:"volume"`
	AdjHigh     float64 `json:"adj_high"`
	AdjLow      float64 `json:"adj_low"`
	AdjClose    float64 `json:"adj_close"`
	AdjOpen     float64 `json:"adj_open"`
	AdjVolume   float64 `json:"adj_volume"`
	Splitfactor float64 `json:"splitfactor"`
	Dividend    float64 `json:"dividend"`
	Symbol      string  `json:"symbol"`
	Exchange    string  `json:"exchange"`
	Date        string  `json:"date"`
}

func NewEODData(logger *zap.Logger) *NewEOD {
	return &NewEOD{
		logger: logger,
	}
}

func (ed *NewEOD) GetEndOfDayData(ctx context.Context, symbol string) (*EOD, error) {
	if err := godotenv.Load(); err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("STOCK_BASE_EOD")+"&symbols="+symbol, nil)
	if err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()
	dataBytes, err := io.ReadAll(res.Body)
	if err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	var val EOD
	if err := json.Unmarshal(dataBytes, &val); err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrEODOutputStream)
	}
	return &val, nil
}

func (ed *NewEOD) LatestEndOfDayData(ctx context.Context, symbol string) (*EOD, error) {
	if err := godotenv.Load(); err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, os.Getenv("STOCK_BASE_EOD_LATEST")+"&symbols="+symbol, nil)
	if err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()
	dataBytes, err := io.ReadAll(res.Body)
	if err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, err
	}
	var val EOD
	if err := json.Unmarshal(dataBytes, &val); err != nil {
		ed.logger.Error(err.Error(), zap.Error(err))
		return nil, errors.New(ErrEODOutputStream)
	}
	return &val, nil
}
