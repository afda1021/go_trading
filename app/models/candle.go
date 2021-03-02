package models

import (
	"fmt"
	"go_trading/bitflyer"
	"time"
)

type Candle struct {
	Duration time.Duration `json:"duration"`
	Time     time.Time     `json:"time"`
	Open     float64       `json:"open"`
	Close    float64       `json:"close"`
	High     float64       `json:"high"`
	Low      float64       `json:"low"`
	Volume   float64       `json:"volume"`
}

func NewCandle(duration time.Duration, timeDate time.Time, open, close, high, low, volume float64) *Candle {
	return &Candle{duration, timeDate, open, close, high, low, volume}
}

/* テーブル名 */
func (c *Candle) TableName() string {
	return fmt.Sprintf("BTC_USD_%s", c.Duration)
}

/* DBに現時刻のcandleを追加 */
func (c *Candle) Create() {
	cmd := fmt.Sprintf("INSERT INTO %s (time, open, close, high, low, volume) VALUES (?,?,?,?,?,?)", c.TableName)
	DbConnection.Exec(cmd, c.Time.Format(time.RFC3339), c.Open, c.Close, c.High, c.Low, c.Volume)
}

func (c *Candle) Save() {
	cmd := fmt.Sprintf("UPDATE %s SET open = ?, close = ?, high = ?, low = ?, volume = ? WHERE time = ?", c.TableName)
	DbConnection.Exec(cmd, c.Open, c.Close, c.High, c.Low, c.Volume, c.Time.Format(time.RFC3339))
}

/* 対応するテーブルから現在時刻と一致するデータを取得(存在しない場合は、nil) */
func GetCandle(duration time.Duration, dateTime time.Time) *Candle {
	tableName := fmt.Sprintf("BTC_USD_%s", duration)
	cmd := fmt.Sprintf("SELECT time, open, close, high, low, volume FROM %s WHERE time = ?", tableName) //WHERE time = ?
	row := DbConnection.QueryRow(cmd, dateTime.Format(time.RFC3339))                                    //dateTime.Format(time.RFC3339)
	var candle Candle
	var times string //Scanはstring型で返す
	err := row.Scan(&times, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
	candle.Time, _ = time.Parse(time.RFC3339, times) //stringからtime.Time型に変換
	if err != nil {
		return nil //一致するデータが存在しない場合
	}
	return NewCandle(duration, candle.Time, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
}

/* tickerが来る度に呼び出し、ticker情報をDBに格納 */
func CreateCandleWithDuration(ticker bitflyer.Ticker, duration time.Duration) {
	currentCandle := GetCandle(duration, ticker.TruncateDateTime(duration)) //対応するテーブルから現在時刻と一致するデータを取得
	price := (ticker.BestBid + ticker.BestAsk) / 2
	if currentCandle == nil { //一致するcandleがない場合、新規作成
		candle := NewCandle(duration, ticker.TruncateDateTime(duration), price, price, price, price, ticker.Volume)
		candle.Create() //DBに現時刻のcandleを追加
		return
	}
	//一致するcandleがある場合、更新
	if currentCandle.High <= price {
		currentCandle.High = price
	} else if currentCandle.Low >= price {
		currentCandle.Low = price
	}
	currentCandle.Volume += ticker.Volume
	currentCandle.Close = price
	currentCandle.Save() //DBのcandleを更新
}
