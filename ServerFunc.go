package main
import "net/http"
import "github.com/gorilla/mux"
import "github.com/go-redis/redis"
import "encoding/json"
import "labix.org/v2/mgo"
import "labix.org/v2/mgo/bson"
import "strconv"
import "fmt"
import "github.com/gorilla/websocket"
import "bytes"
import "io/ioutil"
import "./Redis"
import "./scrapy"
import "./ipsave"
func check(e error){
	if e != nil{
		panic(e)
	}
}
//定义mongoDB数据类型
type Price struct {
	Id bson.ObjectId `bson:"_id"`
	Time int `bson:"time"`
	Close float32 `bson:"close"`
	High float32 `bson:"high"`
	Low float32 `bson:"low"`
	Open float32 `bson:"open"`
	Volumefrom float32 `bson:"volumefrom"`
	Volumeto float32 `bson:"volumeto"`
}
//websocket协议
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
//获取BTC即时行情
func HttpCurrentDataGet(w http.ResponseWriter, r *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr : "localhost:6379",
		Password: "",
		DB: 1,
	})
	defer client.Close()

//获取clien_ip:
	xrealip := ipsave.RemoteIp(r)
	fmt.Println(xrealip)
	allowsecond := ipsave.IpLimitSecond(xrealip)
	if allowsecond==true{
		allowminute := ipsave.IpLimitMinute(xrealip)
		if allowminute==true{
			p,_ :=client.Get("USD").Result()
			data := make(map[string]string)
			data["USD"]=p
			json.NewEncoder(w).Encode(data)
		}else{
			w.Write([]byte("minute down"))
		}
	}else{
		w.Write([]byte("second down"))
	}
}
//websocket 提供btc即时行情
func WebsocketCurrentDataGet(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w,r,nil)
	check(err)
	for {

		client := redis.NewClient(&redis.Options{
			Addr : "localhost:6379",
			Password: "",
			DB:1,
		})
		defer client.Close()
		p,_ := client.Get("USD").Result()
		data := make(map[string]string)
		data["USD"]=p
		msgType, msg, err := conn.ReadMessage()
		if string(msg) == "ping"{
			bytesBuffer := bytes.NewBuffer([]byte{})

			json.NewEncoder(bytesBuffer).Encode(data)
			fmt.Println(bytesBuffer)
			bytestring,_ := ioutil.ReadAll(bytesBuffer)
			fmt.Println(bytestring)
			err = conn.WriteMessage(msgType,bytestring)
			check(err)
			fmt.Println("send",bytestring)

		}else{
			conn.Close()
			fmt.Println(string(msg))
			return
		}
	}
}


//获取btc指定日期行情(美元)
func HttpHistoryDataGet(w http.ResponseWriter, r *http.Request){

//获取clien_ip:
	xrealip := ipsave.RemoteIp(r)
	fmt.Println(xrealip)
	allowsecond := ipsave.IpLimitSecond(xrealip)
	if allowsecond==true{
		allowminute := ipsave.IpLimitMinute(xrealip)
		if allowminute==true{
			params := mux.Vars(r)
			params_int,_ := strconv.ParseInt(params["time"],10,0)
			session, _ := mgo.Dial("127.0.0.1:27017")
			defer session.Close()
			session.SetMode(mgo.Monotonic,true)
			conn := session.DB("history").C("price")
			var price Price
			conn.Find(bson.M{"time":params_int}).One(&price)
			json.NewEncoder(w).Encode(price)
		}else{
			w.Write([]byte("minute down"))
		}

	}else{
		w.Write([]byte("second down"))
}
}
//获取切片数据(时间范围)
func HttpHistoryDataGetSlice(w http.ResponseWriter, r *http.Request){
	xrealip := ipsave.RemoteIp(r)
	fmt.Println(xrealip)
	allowsecond := ipsave.IpLimitSecond(xrealip)
	if allowsecond==true{
		allowminute := ipsave.IpLimitMinute(xrealip)
		if allowminute==true{
			params := mux.Vars(r)
			start_int,_ := strconv.ParseInt(params["startTime"],10,0)
			end_int,_ := strconv.ParseInt(params["endTime"],10,0)
			session,_ := mgo.Dial("127.0.0.1:27017")
			defer session.Close()
			session.SetMode(mgo.Monotonic,true)
			conn := session.DB("history").C("price")
			var price []Price
			conn.Find(bson.M{"time":bson.M{"$gte": start_int,"$lte":end_int}}).All(&price)
			json.NewEncoder(w).Encode(price)
			fmt.Println(start_int,end_int)
		}else{
			w.Write([]byte("minute down"))
		}

	}else{
		w.Write([]byte("second down"))
}


}
//获取btc全部历史数据
func HttpHistoryDataGetAll(w http.ResponseWriter, r *http.Request){
	xrealip := ipsave.RemoteIp(r)
	fmt.Println(xrealip)
	allowsecond := ipsave.IpLimitSecond(xrealip)
	if allowsecond==true{
		allowminute := ipsave.IpLimitMinute(xrealip)
		if allowminute==true{
			session, _ := mgo.Dial("127.0.0.1:27017")
			defer session.Close()
			session.SetMode(mgo.Monotonic,true)
			conn := session.DB("history").C("price")
			var priceAll []Price
			conn.Find(nil).All(&priceAll)
			json.NewEncoder(w).Encode(priceAll)
		}else{
			w.Write([]byte("minute down"))
		}

	}else{
		w.Write([]byte("second down"))
}
}





//获取全部货币即时行情
func WebsocketAllCurrentDataGet(w http.ResponseWriter, r *http.Request){

	conn ,err := upgrader.Upgrade(w,r,nil)
	check(err)
	for {

		msgType, msg, err := conn.ReadMessage()
		if string(msg) == "ping"{

			datass := test.ScrapyDataGet()
			bytesBuffer := bytes.NewBuffer([]byte{})
			json.NewEncoder(bytesBuffer).Encode(datass)
			mybyte,_ :=ioutil.ReadAll(bytesBuffer)
			err = conn.WriteMessage(msgType,mybyte)
			check(err)
		}else{
			conn.Close()
			fmt.Println(string(msg))
			return
		}
	}
}
//获取全部货币即时行情http
func HttpAllCurrentDataGet(w http.ResponseWriter, r *http.Request){

	xrealip := ipsave.RemoteIp(r)
	fmt.Println(xrealip)
	allowsecond := ipsave.IpLimitSecond(xrealip)
	if allowsecond==true{
		allowminute := ipsave.IpLimitMinute(xrealip)
		if allowminute==true{
			FinalData :=test.ScrapyDataGet()
			json.NewEncoder(w).Encode(FinalData)
		}else{
			w.Write([]byte("minute down"))
		}

	}else{
		w.Write([]byte("second down"))
}

}

func main(){
	go CurrentDataGet()
	go HistoryDataGet()
	go myscrapy.ScrapyData()
	router := mux.NewRouter()
	router.HandleFunc("/current_data",HttpCurrentDataGet).Methods("GET")
	router.HandleFunc("/history_data/{time}",HttpHistoryDataGet).Methods("GET")
	router.HandleFunc("/history_data",HttpHistoryDataGetAll).Methods("GET")
	router.HandleFunc("/history_data/{startTime}/{endTime}",HttpHistoryDataGetSlice).Methods("GET")
	router.HandleFunc("/current_data/websocket",WebsocketCurrentDataGet).Methods("GET")
	router.HandleFunc("/allcoin_current",HttpAllCurrentDataGet)
	router.HandleFunc("/allcoin_current/websocket",WebsocketAllCurrentDataGet)
	http.ListenAndServe("127.0.0.1:8080",router)
}
