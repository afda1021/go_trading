package controllers

import (
	"fmt"
	"go_trading/app/models"
	"go_trading/bitflyer"
)

func StreamIngestionData() {
	var tickerChannel = make(chan bitflyer.Ticker) //ticker用チャネルを生成
	go bitflyer.GetRealTimeTicker(tickerChannel)   //tickerChannelに取得したtickerを入れていく

	for ticker := range tickerChannel {
		fmt.Println(ticker)
		fmt.Println(models.DbConnection) //DB接続
	}
}
