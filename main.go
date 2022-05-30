package main

import (
	//"io/ioutil"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	//"time"

	"github.com/joho/godotenv"
)
var protocol string
var target string
var port string
var tps int
var rpt int

func initenv() {
	var e string
	godotenv.Load(".env")
	if e = os.Getenv("PROTOCOL"); e != "" {
		protocol = e
	} else {
		protocol = "https"
	}
	if e = os.Getenv("TARGET"); e != "" {
		target = e
	} else {
		panic("no target, pussy")
	}
	if e = os.Getenv("PORT"); e != "" {
		port = e
	} else {
		port = "443"
	}
	// if e = os.Getenv("USER"); e != "" {
	// 	user = e
	// } else {
	// 	panic("no user")
	// }
	// if e = os.Getenv("PASS"); e != "" {
	// 	pass = e
	// } else {
	// 	panic("no pass")
	// }
	if e = os.Getenv("THREADS_PER_SECOND"); e != "" {
		tps, _ = strconv.Atoi(e)
	} else {
		panic("no TPS")
	}
	if e = os.Getenv("REQUESTS_PER_THREAD"); e != "" {
		rpt, _ = strconv.Atoi(e)
	} else {
		panic("no RPT")
	}
	

}

func main() {

	initenv()
	InitClear()


	client := &http.Client{}
	targetUrl := genTarget(protocol, target, port)

	statuscode := make(chan int)
	successfulcalls := 0
	failedcalls := 0
	go countCalls(statuscode, &successfulcalls, &failedcalls)

	secondticker := time.NewTicker(time.Second)
	for range secondticker.C {		
		for i := 0; i < tps; i++ {
			go call(client, "GET", targetUrl, statuscode)
		}
		printInfos(successfulcalls, failedcalls)
	}



}

func call(client *http.Client, method string, url string, statuscode chan int) {
	
	req, _ := http.NewRequest(method, url, nil)
	//req.SetBasicAuth(user, pass)

	for i := 1; i <= rpt; i++ {
		resp, err := client.Do(req)
		if err != nil {
			statuscode <- 500
			return
		}
		statuscode <- resp.StatusCode
	}
}

func countCalls(c chan int, success *int, failure *int ) {
	for i := range c {
		if i == 200 {
			*success +=1
		} else  {
			*failure +=1
		} 
	}
}

func genTarget(protocol string, target string, port string) string {
	return protocol + "://" + target + ":" + port
}
func printInfos(successful int, failed int) {
	CallClear()
	fmt.Printf("successful: %d\n", successful)
	fmt.Printf("failed: %d\n", failed)
}

//clear terminal
