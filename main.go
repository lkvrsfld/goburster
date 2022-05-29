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

var target string
// var user string
// var pass string
var rps int

var failedcalls int
var successfulcalls int

func initenv() {
	var e string
	godotenv.Load(".env")
	if e = os.Getenv("TARGET"); e != "" {
		target = e
	} else {
		panic("no target, pussy")
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
	if e = os.Getenv("REQUESTS_PER_SECOND"); e != "" {
		rps, _ = strconv.Atoi(e)
	} else {
		panic("no RPS")
	}

}

func main() {

	initenv()
	InitClear()


	client := &http.Client{}

	statuscode := make(chan int)

	successfulcalls := 0
	failedcalls := 0

	go countCalls(statuscode, &successfulcalls, &failedcalls)
	secondticker := time.NewTicker(time.Second)
	for range secondticker.C {
		CallClear()
		for i := 0; i < rps; i++ {
			go call(client, "GET", "https://" + target, statuscode)
		}
		printInfos(successfulcalls, failedcalls)
	}



}

func call(client *http.Client, method string, url string, statuscode chan int) {
	
	req, _ := http.NewRequest(method, url, nil)
	//req.SetBasicAuth(user, pass)

	for i := 1; i <= 10; i++ {
		resp, _ := client.Do(req)
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

func printInfos(successful int, failed int) {
	fmt.Printf("successful: %d\n", successful)
	fmt.Printf("failed: %d\n", failed)
}

//clear terminal
