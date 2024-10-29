package handler

import (
	"html/template"
	"log"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	pageTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Handler, Root, Error accessing HTML file: %v", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Unable to find page."))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err = pageTemplate.Execute(w, nil)
	if err != nil {
		//Too late to do any real error handling, just log the error.
		log.Printf("Handler, Root, Error executing page: %v", err)
	}
}
