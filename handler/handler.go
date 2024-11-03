package handler

import (
	//"fmt"
	"html/template"
	"log"
	"net/http"

	//"strings"
	"github.com/JohnDirewolf/capstone/maze"
	"github.com/JohnDirewolf/capstone/shared/types"
)

// STRUCTURES
type pageData struct {
	Title       string
	Rooms       template.HTML
	Compass     template.HTML
	Description string
}

// HANDLERS
func Root(w http.ResponseWriter, r *http.Request) {
	pageData := struct {
		Title string
	}{
		Title: "Root of Capstone Project!",
	}

	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/index.html")
	if err != nil {
		log.Printf("Handler, Root, Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err = pageTemplate.Execute(w, pageData)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, Root, Error executing page: %v", err)
	}
}

func Game(w http.ResponseWriter, r *http.Request) {
	//This takes the action encoded in r and updates and displays the new updated page.
	//To prevent caching, add the following header to the response "Cache-Control: no-store or Cache-Control: no-cache, no-store, must-revalidate"

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	//Get the action path parameter from the request.
	action := r.URL.Query().Get("action")
	switch action {
	case "start":
		start(w)
	case "end":
		end(w)
	case types.North, types.South, types.West, types.East:
		move(action, w)
	default:
		//Unknown action
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		var pageBody string = "400: Unknown action."
		w.Write([]byte(pageBody))
	}
}

// ACTIONS
func start(w http.ResponseWriter) {
	//Initialize the Maze
	maze.Init()

	compassHTML, roomTitle, roomDescription := maze.GetRoomInfo()
	startPageData := pageData{
		Title:       roomTitle + "- Start",
		Rooms:       maze.GenerateKnownMap(),
		Compass:     compassHTML,
		Description: "You enter the maze, the entrance slamming shut behind you. " + roomDescription,
	}

	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/maze.html")
	if err != nil {
		log.Printf("Handler, Game, start Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	//log.Printf("Known Map: %s\n", startPageData.Rooms)
	err = pageTemplate.Execute(w, startPageData)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, Game, start, Error executing page: %v", err)
	}
}

func end(w http.ResponseWriter) {
	startPageData := pageData{
		Title: "Maze Runner End!",
	}

	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/end.html")
	if err != nil {
		log.Printf("Handler, Game, end, Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	err = pageTemplate.Execute(w, startPageData)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, Game, end, Error executing page: %v", err)
	}
}

func move(direction string, w http.ResponseWriter) {
	maze.Move(direction)

	compassHTML, roomTitle, roomDescription := maze.GetRoomInfo()
	startPageData := pageData{
		Title:       roomTitle,
		Rooms:       maze.GenerateKnownMap(),
		Compass:     compassHTML,
		Description: roomDescription,
	}

	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/maze.html")
	if err != nil {
		log.Printf("Handler, Game, move, Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	err = pageTemplate.Execute(w, startPageData)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, Game, move, Error executing page: %v", err)
	}
}
