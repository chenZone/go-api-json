package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"time"
	"github.com/go-redis/redis"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)



func HistoryDataGet(){
	session, _:=mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic,true)
	conn := session.DB("history").C("price")
	resp, _ :=http.Get("https://min-api.cryptocompare.com/data/histoday?fsym=BTC&tsym=USD&allData=true")
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	var m map[string]interface{}
	err := json.Unmarshal(body,&m)
	conn.RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}else{
		strs := m["Data"].([]interface{})
		for _,value := range strs{
/*			change_value := value.(map[string]interface{})

			time_str := change_value["time"].(float64)
			time_int := int(time_str)
			real_time := time.Unix(int64(time_int),0)
			close_str := change_value["close"].(float64)
			close_int := int(close_str)
			high_str := change_value["high"].(float64)
			high_int := int(high_str)
			low_str := change_value["low"].(float64)
			low_int := int(low_str)
			open_str := change_value["open"].(float64)
			open_int := int(open_str)
			volumefrom_str := change_value["volumefrom"].(float64)
			volumefrom_int := int(volumefrom_str)
			volumeto_str := change_value["volumeto"].(float64)
			volumeto_int := int(volumeto_str)
			
			fmt.Println(close_int)
			fmt.Println(time_int)
			fmt.Println(high_int)
			fmt.Println(low_int)
			fmt.Println(open_int)
			fmt.Println(volumefrom_int)
			fmt.Println(volumeto_int)
			fmt.Println(real_time)
			enc := json.NewEncoder(os.Stdout)
			enc.Encode(change_value)
		}*/
//		env := json.NewEncoder(os.Stdout)
//		env.Encode(m)
		conn.Insert(value)
	}
		t:=time.NewTimer(24*time.Hour)
		<-t.C
}
}

func CurrentDataGet(){
	client := redis.NewClient(&redis.Options{
		Addr : "localhost:6379",
		Password: "",
		DB : 1,
	})
	defer client.Close()
	for{
		t := time.NewTimer(1 * time.Second)
		<-t.C
		resp, _ := http.Get("https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD")
		defer resp.Body.Close()
		body,_ := ioutil.ReadAll(resp.Body)
		var m map[string]float64
		json.Unmarshal(body,&m)
		client.Set("USD",m["USD"],0)
		fmt.Println("update!")
	}
}





