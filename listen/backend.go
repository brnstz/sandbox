package main

import (
	"log"
	"os"
)

var debug bool

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if os.Getenv("LISTEN_DEBUG") == "1" {
		debug = true
	}

	for i := 0; i < numWorkers; i++ {
		var s show
		var b band
		var v venue
		go recvAndWrite(&s)
		go recvAndWrite(&b)
		go recvAndWrite(&v)
	}

	go pullFromOMR()
	go agg()

	//http.HandleFunc("/api/shows.json", getShows)

	// By default look for a static asset
	//http.Handle("/", http.FileServer(http.Dir(staticDir)))

	/*
		err := http.ListenAndServe(":8003", nil)
		if err != nil {
			log.Fatal(err)
		}
	*/

	forever := make(chan bool)
	<-forever
}
