package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal([]string{"Hello", "World"})
		if err != nil {
			fmt.Fprintf(w, "Error %v", err)
			return
		}
		fmt.Fprintf(w, "%v", string(b))
	}
	r := http.NewServeMux()
	r.HandleFunc("/bar", handler)
	srv := &http.Server{
		Handler: r,
		Addr:    ":5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() { log.Fatal(srv.ListenAndServe()) }()

	r2 := http.NewServeMux()
	r2.HandleFunc("/foo", handler)
	srv2 := &http.Server{
		Handler: r2,
		Addr:    ":5001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("http://localhost:5000/bar")
	fmt.Println("http://localhost:5001/foo")
	fmt.Println("http http://localhost:4000/bar && http http://localhost:4001/foo")
	log.Fatal(srv2.ListenAndServe())
}
