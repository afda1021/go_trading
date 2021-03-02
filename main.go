package main

import (
	"go_trading/app/controllers"
)

func main() {
	// fmt.Println(models.DbConnection) //DB接続
	controllers.StreamIngestionData() //ticker情報をDBに格納
}
