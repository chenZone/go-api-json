package main


func main() {
	done := make(chan bool)
	scrapy_data()
	Get_data()
	Http_server(done)
	<-done
}
