package bitflyer

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

/* JSON-RPCのプロトコルを定義 */
type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result.omitempty"`
	Id      *int        `json:"id.omitempty"`
}

/* bitflyerのJSON-RPCプロトコルが、名前付き引数("channel")での利用を想定して設計されているため、別途typeを定義 */
type SubscribeParams struct {
	Channel string `json:"channel"`
}

//tickerChannelに取得したtickerを渡す
//JsonRPC2を送信(WriteJSON)するとJsonRPC2のParamsにmessageが入って返って(ReadJSON)くる
//JsonRPC2のMethodは、送る時がsubscribeで、返ってくる時がchannelMessage
func GetRealTimeTicker(ch chan<- Ticker) {
	u := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "json-rpc"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil) //websocketによる接続を確立
	if err != nil {
		log.Fatal("websocket:", err)
		return
	}
	defer c.Close()
	err = c.WriteJSON(&JsonRPC2{Version: "2.0", Method: "subscribe", Params: &SubscribeParams{Channel: "lightning_ticker_BTC_USD"}}) //JSON形式で送信。subscribe：チャンネルの購読を開始
	if err != nil {
		log.Fatal("subscribe:", err)
		return
	}

OUTER:
	for {
		message := new(JsonRPC2)   //JsonRPC2のポインタ型を生成、"message := &JsonRPC2{}"でも良い
		err := c.ReadJSON(message) //JSON形式で受信。プロトコルは、定義したJsonRPC2
		if err != nil {
			log.Println("read:", err)
			return
		}
		if message.Method == "channelMessage" {
			switch v := message.Params.(type) { //v = map[channel:lightning_ticker_BTC_USD message:map[(ticker)...]]
			case map[string]interface{}:
				for key, binary := range v { //key = [channel, message]
					if key == "message" {
						marshaTic, err := json.Marshal(binary) //構造体→json
						if err != nil {
							continue OUTER
						}
						var ticker Ticker
						if err := json.Unmarshal(marshaTic, &ticker); err != nil {
							continue OUTER
						}
						ch <- ticker //チャネルにtickerを渡す
					}
				}
			}
		}
	}
}
