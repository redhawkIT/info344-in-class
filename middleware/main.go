package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := "localhost:4000"

	mux := http.NewServeMux()
	muxLogged := http.NewServeMux()
	muxLogged.HandleFunc("/v1/hello1", HelloHandler1)
	muxLogged.HandleFunc("/v1/hello2", HelloHandler2)
	mux.HandleFunc("/v1/hello3", HelloHandler3)

	mux.Handle("/v1/", logRequests(muxLogged))

	// http.HandleFunc("/v1/hello1/", logReqs(HelloHandler1))
	// http.HandleFunc("/v1/hello2/", logReqs(HelloHandler2))
	// http.HandleFunc("/v1/hello3/", logReqs(HelloHandler3))

	// log.Fata

	logger := log.New(os.Stout, "", log.LstdFlags)
	mux.Handle("/v1/", logRequests(logger)(muxLogged))

	fmt.Printf("listening at %s...\n", addr)
	// log.Fatal(http.ListenAndServe(addr, nil))
	log.Fatal(http.ListenAndServe(addr, mux))
}
