package models

import (
	"time"

	"github.com/markcheno/go-talib"
)

/* 読み取ったcandleの格納用 */
type DataFrameCandle struct {
	Duration time.Duration `json:"duration"`
	Candles  []Candle      `json:"candles"`
	Smas     []Sma         `json:"smas,omitempty"`
}

type Sma struct {
	Period int       `json:"period,omitempty"` //omitempty:0とか値がない場合、jsonにしたときに省略
	Values []float64 `json:"values,omitempty"`
}

/* Closeの配列 */
func (df *DataFrameCandle) Closes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Close
	}
	return s
}

/* periodに対するSMAを計算し、df.Smasに追加 */
func (df *DataFrameCandle) AddSma(period int) bool {
	if len(df.Candles) > period { //period以上データがないと作れない
		df.Smas = append(df.Smas, Sma{
			Period: period,
			Values: talib.Sma(df.Closes(), period),
		})
		return true
	}
	return false //periodがデータ数を超えている場合は計算不能
}
