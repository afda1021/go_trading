package controllers

import (
	"encoding/json"
	"go_trading/app/models"
	"go_trading/config"
	"net/http"
	"strconv"
	"text/template"
)

var tempaltes = template.Must(template.ParseFiles("app/views/chart.html"))

/* Google Candle Stick Chartsを表示 */
func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	err := tempaltes.ExecuteTemplate(w, "chart.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* Ajaxでクライアントのparamsを取得しjsonを送信 */
func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if strLimit == "" || err != nil || limit < 0 || limit > 1000 {
		limit = 1000
	}

	duration := r.URL.Query().Get("duration")
	if duration == "" {
		duration = "1m"
	}
	durationTime := config.Durations[duration]

	df, _ := models.GetAllCandle(durationTime, limit) //指定されたテーブルから全てのcandleを取得

	js, err := json.Marshal(df) //構造体→json
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js) //jsonにしたものを返す
}

func StartWebServer() {
	http.HandleFunc("/api/candle/", apiCandleHandler) //Ajaxでクライアントのparamsを元にjsonを送信
	http.HandleFunc("/chart/", viewChartHandler)      //chart.htmlへ遷移
	http.ListenAndServe(":8080", nil)
}
