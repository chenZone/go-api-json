package main

import "os/exec"
import "fmt"
import "time"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func scrapy_data() {

		start := time.Now()
		cmd := exec.Command("bash", "-c", "cd /root/feixiaohao-spider/feixiaohao && scrapy crawl feixiao")
		cmd.Start()
		fmt.Println("开始抓取")
		cmd.Wait()
		fmt.Println("finished!")
		end := time.Now()
		delta := end.Sub(start)
		fmt.Printf("take time:%s\n", delta)
		timer1 := time.NewTimer(time.Second * 2)
		fmt.Printf("waiting..")
		<-timer1.C
		fmt.Printf("done!\n")

}
