package main

import (
	"fmt"
	"log"
	"net/http"
)

func main()  {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("aaaaa")
	})
	err := http.ListenAndServe("0.0.0.0", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
