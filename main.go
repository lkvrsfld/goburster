package main

import (
	//"io/ioutil"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"time"

	"github.com/joho/godotenv"
)

var target string
var user string
var pass string
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
	if e = os.Getenv("USER"); e != "" {
		user = e
	} else {
		panic("no user")
	}
	if e = os.Getenv("PASS"); e != "" {
		pass = e
	} else {
		panic("no pass")
	}
	if e = os.Getenv("REQUESTS_PER_SECOND"); e != "" {
		rps, _ = strconv.Atoi(e)
	} else {
		panic("no RPS")
	}

}

func main() {

	initenv()
	InitClear()

	// here we want to imput sth like 100req/s

	secondticker := time.NewTicker(time.Second)
	for range secondticker.C {

		for i := 0; i < rps; i++ {
			go call("GET", "https://" + target)
		}
		printInfos()
	}



}

func call(method string, url string) {
	client := &http.Client{}

	req, _ := http.NewRequest(method, url, nil)
	//req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	fmt.Println(resp)

	if resp.StatusCode == 200 && err == nil {
		successfulcalls += 1
	} else {
		failedcalls += 1
	}

	/* fmt.Println(resp.StatusCode) */
}

func printInfos() {
	fmt.Printf("successful: %d\n", successfulcalls)
	fmt.Printf("failed: %d\n", failedcalls)
}

//clear terminal
