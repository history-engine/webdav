package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	port     string
	data     string
	username string
	password string
)

func init() {
	flag.StringVar(&data, "data", "", "set data directory")
	flag.StringVar(&port, "port", "2233", "set listen port")
	flag.StringVar(&username, "username", "admin", "set username")
	flag.StringVar(&password, "password", "admin", "set password")
}

func main() {
	flag.Parse()

	data = checkData(data)
	if data == "" || port == "" || username == "" || password == "" {
		flag.Usage()
		return
	}

	log.Printf("data directory: %s\n", data)
	log.Printf("webdav serve at: %s\n", port)

	http.HandleFunc("/", indexHandler)
	if err := http.ListenAndServe(":"+port, nil); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe %s err %v\n", port, err)
	}
}
