package test


import "github.com/go-redis/redis"
import "fmt"
import "time"
import "sort"
import "strconv"
func DataStorage (a,b string) {
	client := redis.NewClient(&redis.Options{
		Addr : "localhost:6379",
		Password: "",
		DB: 1,
	})
	defer client.Close()
	client.Set(a,b,0)

}

func DataGet (a string) {
	client := redis.NewClient(&redis.Options{
		Addr : "localhost:6379",
		Password: "",
		DB: 1,
	})

	defer client.Close()
	for {

		p,_ :=client.Get(a).Result()
		fmt.Println(p)
		t := time.NewTimer(2 * time.Second)
		<-t.C
	}
}


type Data struct {
	Name string `json:"name"`
	F_id int `json:"f_id"`
	Flow_market_price string `json:"flow_market_price"`
	Price string `json:"price"`
	Flow_amount string `json:Flow_amount`
	Trade_amount string `json:Trade_amount`
	Price_change string `json:Price_change`
}
type Datas []Data
var datas = Datas{}

func (d Datas) Swap(i,j int) {
	d[i],d[j] = d[j],d[i]
}

func (d Datas) Len() int {
	return len(d)
}

func (d Datas) Less(i,j int) bool {
	return d[i].F_id < d[j].F_id
}

func RedisGetAll(st ...string) []Data {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB:0,
	})
	defer client.Close()
	datas = datas[:0:0]
	for _,i := range(st){
		hash,_ :=client.HGetAll(i).Result()
		hash_int, _ := strconv.Atoi(hash["f_id"])
		data := Data{
			Name:    hash["name"],
			F_id:    hash_int,
			Flow_market_price: hash["flow_market_price"],
			Price: hash["price"],
			Flow_amount: hash["flow_amount"],
			Trade_amount: hash["trade_amount"],
			Price_change: hash["price_change"],
		}
		datas = append(datas,data)
	}
	sort.Stable(datas)
	return datas
}
func ScrapyDataGet() []Data{
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB:0,
	})
	defer client.Close()
	result, _ :=client.Keys("*").Result()
	FinalData := RedisGetAll(result...)
	return FinalData
}



