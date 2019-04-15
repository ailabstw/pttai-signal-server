package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")

	flag.Parse()
	log.SetFlags(0)

	server := NewServer()

	http.HandleFunc("/signal", server.SignalHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
