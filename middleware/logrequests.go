package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// func logReq(r *http.Request) {
// 	log.Println(r.Method, r.URL.Path)
// }

func logReqs(hfn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		start := time.Now()
		hfn(w, r) //Execute function
		fmt.Printf("%v\n", time.Since(start))
	}
}

func logRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("%v\n", time.Since(start))
	})
}
