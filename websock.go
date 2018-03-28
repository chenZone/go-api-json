package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"encoding/json"
	"net/http"
	"sort"
	"log"
	"github.com/gorilla/websocket"
	"strconv"
)

type Data struct {
	Name               string       `json:"name"`
	F_id               int          `json:"f_id"`
	Flow_market_price  string       `json:"flow_market_price"`
	Price              string       `json:"price"`
	Flow_amount        string       `json:"flow_amount"`
	Trade_amount       string       `json:"trade_amount"`
	Price_change       string       `json:"price_change"`
}

type Datas []Data
var count = 1
var datas = Datas{}
func (d Datas) Len() int {
	return len(d)
}

func (d Datas) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Datas) Less(i, j int) bool {
	return d[i].F_id < d[j].F_id
}
func Get_data() {
	datas = datas[:0:0]
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	result, _ := client.Keys("*").Result()
	for _, i := range result {
		hash, _ := client.HGetAll(i).Result()
		hash_int, _ := strconv.Atoi(hash["f_id"])
		data := Data{
			Name:              hash["name"],
			F_id:              hash_int,
			Flow_market_price: hash["flow_market_price"],
			Price:             hash["price"],
			Flow_amount:       hash["flow_amount"],
			Trade_amount:      hash["trade_amount"],
			Price_change:      hash["price_change"],
		}
		datas = append(datas,data)

	}
	sort.Stable(datas)
	fmt.Println("redis data get ..finished")
}


func check(e error){
	if e != nil{
		panic(e)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/",home)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))

}

func home(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	check(err)
	Get_data()
	for {
		msgType, msg, err := conn.ReadMessage()
		check(err)
		if string(msg) == "ping" {
			if count != 1{
				Get_data()
				datass := Return_json()
				fmt.Println("first_received")
				err = conn.WriteMessage(msgType,datass)
				check(err)
			}else{

			
			Scrapy_data()
			Get_data()
			datass := Return_json()
			fmt.Println("received")
			err = conn.WriteMessage(msgType,datass)
			check(err)
			}
		} else {
			conn.Close()
			fmt.Println(string(msg))
			return
		}
	}
}




func Return_json () []byte {
	strsum := []byte("")
	str1 := []byte("[")
	str2 := []byte(",")
	str3 := []byte("]")
	strsum = append(strsum,str1...)
	for key,value := range datas {
		vr,_ := json.Marshal(value)
		strsum = append(strsum,vr...)
		if key!=len(datas)-1{
			strsum = append(strsum,str2...)
		}

	}
	strsum = append(strsum,str3...)
//	last_data := string(strsum)
	return strsum
}

