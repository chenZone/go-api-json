package ipsave
import "time"
import "fmt"
import "github.com/go-redis/redis"
import "strconv"
func IpLimitMinute(ip string) bool {
	client := redis.NewClient(&redis.Options{
		Addr : "localhost:6379",
		Password: "",
		DB : 2,
	})
	defer client.Close()
	defer func(){
		if x:=recover();x != nil{
			fmt.Println("超速")
		}
	}()
	if _,err:= client.Get(ip).Result();err==nil {
		p,_ :=client.Get(ip+"count").Result()
		pint,_ := strconv.ParseInt(p,10,0)
		pint -= 1
		client.Set(ip+"count",pint,0)
		if pint <= 0{
			return false
		}else{
			return true
		}
	}else{
		client.Set(ip,"exist",time.Second*60)
		client.Set(ip+"count",60,0)

	}
	return true
}



func IpLimitSecond(ip string) bool {
	client := redis.NewClient(&redis.Options{
		Addr : "localhost:6379",
		Password: "",
		DB : 3,
	})
	defer client.Close()
	defer func(){
		if x:=recover();x != nil{
			fmt.Println("超速")
		}
	}()
	if _,err:= client.Get(ip).Result();err==nil {
		p,_ :=client.Get(ip+"count").Result()
		pint,_ := strconv.ParseInt(p,10,0)
		pint -= 1
		client.Set(ip+"count",pint,0)
		if pint <= 0{
			return false
		}else{
			return true
		}
	}else{
		client.Set(ip,"exist",time.Second)
		client.Set(ip+"count",1,0)

	}
	return true
}
