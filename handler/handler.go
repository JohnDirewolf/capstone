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
	//This takes the action in the URL in r and updates and displays the new updated page.
	//To prevent caching, add the following header to the response "Cache-Control: no-store or Cache-Control: no-cache, no-store, must-revalidate"

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	//Get the action path parameter from the request.
	action := types.UrlAction(r.URL.Query().Get("action"))
	switch types.UrlAction(action) {
	case types.Start:
		start(w)
	case types.End:
		end(w)
	case types.North, types.South, types.West, types.East:
		move(action, w)
	case types.Get:
		get(w)
	case types.Use:
		use(w)
	case types.Attack:
		attack(w)
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
	generateRoom(w, types.SpecialStatus{IsStart: true})
}

func end(w http.ResponseWriter) {
	startPageData := types.PageData{
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

func move(direction types.UrlAction, w http.ResponseWriter) {
	generateRoom(w, maze.Move(direction))
}

func get(w http.ResponseWriter) {
	maze.GetItems()
	generateRoom(w, types.SpecialStatus{})
}

func use(w http.ResponseWriter) {
	maze.UseKey()
	generateRoom(w, types.SpecialStatus{Unlocked: true})
}

func attack(w http.ResponseWriter) {
	//This checks if the player was successful, if so we continue with the creature dead, otherwise we have a loss screen.
	if maze.Attack() {
		generateRoom(w, types.SpecialStatus{Vanquished: true})
		return
	}
	failure(w)
}

func generateRoom(w http.ResponseWriter, special types.SpecialStatus) {
	var pageInfo types.PageData
	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/maze.html")
	if err != nil {
		log.Printf("Handler, generateRoom, Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	pageInfo = maze.GetPageInfo(special)

	err = pageTemplate.Execute(w, pageInfo)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, generateRoom, Error executing page: %v", err)
	}
}

func failure(w http.ResponseWriter) {
	pageTemplate, err := template.ParseFiles("templates/shared/base.html", "templates/shared/header.html", "templates/failure.html")
	if err != nil {
		log.Printf("Handler, failure, Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	pageInfo := types.PageData{
		Title: "You have died!",
	}

	err = pageTemplate.Execute(w, pageInfo)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, failure, Error executing page: %v", err)
	}
}
