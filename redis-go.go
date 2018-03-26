package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"encoding/json"
	"net/http"
	"sort"
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


func Return_json (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	str1 := []byte("[\r")
	str2 := []byte(",")
	str3 := []byte("]")
	w.Write(str1)
	for key,value := range datas {
		vr,_ := json.Marshal(value)
		w.Write(vr)
		if key!=len(datas)-1{
			w.Write(str2)
		}

	}
	w.Write(str3)
}

