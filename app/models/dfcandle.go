package models

import "time"

type DataFrameCandle struct {
	Duration time.Duration `json:"duration"`
	Candles  []Candle      `json:"candles"`
}
