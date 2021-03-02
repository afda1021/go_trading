package controllers

import (
	"fmt"
	"go_trading/app/models"
	"go_trading/bitflyer"
	"go_trading/config"
)

/* ticker情報を取得し、DBに格納 */
func StreamIngestionData() {
	var tickerChannel = make(chan bitflyer.Ticker) //ticker用チャネルを生成
	go bitflyer.GetRealTimeTicker(tickerChannel)   //tickerChannelに取得したtickerを入れていく

	go func() {
		for ticker := range tickerChannel { //tickerChannelにデータが入るたびにfor文が実行される
			for _, duration := range config.Durations { //3つのDBに順番に格納
				models.CreateCandleWithDuration(ticker, duration) //ticker情報をDBに格納
			}
			fmt.Println(ticker)
		}
	}()
}
