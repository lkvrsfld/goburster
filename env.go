package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)


func (b *burster)initFromEnv() (err error) {
	var e string
	err = godotenv.Load(".env")
	// get rhost
	var rhost rhost
	if e = os.Getenv("PROTOCOL"); e != "" {
		rhost.protocol = e
	} else {
		rhost.protocol = "http"
	}
	if e = os.Getenv("TARGET"); e != "" {
		rhost.target = e
	} else {
		err = errors.New("No target specified. please specify")
	}
	if e = os.Getenv("PORT"); e != "" {
		rhost.port = e
	} else {
		rhost.port = "80"
	}

	b.rhost = rhost

	if e = os.Getenv("THREADS_PER_SECOND"); e != "" {
		b.tps, _ = strconv.Atoi(e)
	} else {
		b.tps = 100
		log("no threads per second given, default is 100")
	}
	if e = os.Getenv("REQUESTS_PER_THREAD"); e != "" {
		b.rpt, _ = strconv.Atoi(e)
	} else {
		b.rpt = 100
		log("no requests per thread given, default is 10")
	}
	
	return
}