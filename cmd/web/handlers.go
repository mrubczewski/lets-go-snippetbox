package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Add("Server", "Go")
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	_, err = fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.logger.Error(err.Error())
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Display a form for creating a new snippet..."))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.logger.Error(err.Error())
	}
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("Save a new snippet..."))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.logger.Error(err.Error())
	}
}
