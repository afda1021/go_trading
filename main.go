package main

import (
	"go_trading/app/controllers"
)

func main() {
	controllers.StreamIngestionData() //ticker情報をDBに格納
	controllers.StartWebServer()      //サーバー起動、candleをブラウザに表示
}
