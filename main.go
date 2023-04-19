package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type burster struct {
	client          *http.Client // this will get []*http.Client, when tested, and if it brigns performance improvements
	responseChannel chan int
	results         map[int]int
	rhost           rhost
	timeout         time.Duration
	tps             int // threads per secondticker
	rpt             int // requests per thread

}

type rhost struct {
	protocol string
	target   string
	port     string // why do we need two type casts
}

func main() {
	var burster burster

	//dependencies
	InitClear()

	if err := burster.initFromEnv(); err != nil {
		log("initialisation from .env failed. ")
		panic(err)
	}

	targetUrl := genTarget(burster.rhost.protocol, burster.rhost.target, burster.rhost.port)
	burster.init()

	go burster.countCalls()

	secondticker := time.NewTicker(time.Second)

	for range secondticker.C {
		for i := 0; i < burster.rpt; i++ {
			go burster.call(burster.client, "GET", targetUrl, burster.responseChannel)
		}
		burster.printInfos()
	}
}

func (b *burster) init() (err error) {
	b.client = &http.Client{}
	b.client.Timeout = time.Duration(2 * time.Second)
	b.responseChannel = make(chan int)
	b.results = make(map[int]int)
	return
}

func (b burster) call(client *http.Client, method string, url string, responseChannel chan int) {

	req, _ := http.NewRequest(method, url, nil)

	for i := 1; i <= b.rpt; i++ {
		resp, err := client.Do(req)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			responseChannel <- 999 // means timeout
			return
		} else if err != nil {
			responseChannel <- 666 // means some other error
			return
		}

		responseChannel <- resp.StatusCode
	}
}

func (b *burster) countCalls() {
	for i := range b.responseChannel {
		b.results[i] += 1
	}
}

func (b burster) printInfos() {
	frame := fmt.Sprintf("protocol: %s\ntarget: %s\nport: %s\nrequests per second: %drq/s\nthreads per second: %d\nrequests per thread: %d\nrequests per second: %drq/s\nport: %s\nresponses: \n\n",
	b.rhost.protocol, b.rhost.target, b.rhost.port, b.tps*b.rpt, b.tps, b.rpt, b.tps*b.rpt, b.rhost.port)
	CallClear()

	fmt.Print(frame)

	for statuscode, count := range b.results {
		fmt.Printf("%d: %d\n", statuscode, count)
	}
}

func genTarget(protocol string, target string, port string) string {
	return protocol + "://" + target + ":" + port
}

func log(s string) {
	fmt.Printf("Logger: %s\n", s)
}
