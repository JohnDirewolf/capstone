package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JohnDirewolf/capstone/handler"
)

func main() {
	//Currently the address is const, later will be environmental variable.
	const addr = ":8080"

	//Creating new servemux "dm" which is a nod to Dungeon Master.
	dm := http.NewServeMux()

	//Set up our server
	srv := http.Server{
		Handler:      dm,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	//Set up handlers
	dm.HandleFunc("/", handler.Page)

	fmt.Println("server started on ", addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
