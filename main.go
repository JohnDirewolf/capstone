package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JohnDirewolf/capstone/database"
	"github.com/JohnDirewolf/capstone/handler"
)

func main() {
	//set up database.
	err := database.Init()
	if err != nil {
		log.Println("Error reported from Database.")
		return
	}
	log.Println("Database connected")

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

	//Basic file server for css, images, etc.
	dm.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("./static/css"))))
	dm.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("./assets/images"))))

	//Set up page handlers
	dm.HandleFunc("/", handler.Root)
	dm.HandleFunc("/app", handler.Game)

	fmt.Println("server started on ", addr)
	err = srv.ListenAndServe()
	log.Fatal(err)

	err = database.Close()
	log.Fatal(err)

	//comment
}
