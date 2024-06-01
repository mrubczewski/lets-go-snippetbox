package main

import (
	"errors"
	"fmt"
	"github.com/mrubczewski/lets-go-snippetbox/internal/models"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type snippetCreateForm struct {
	Title       string
	Content     string
	ExpiresAt   int
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippets = snippets
	app.render(w, r, http.StatusOK, "home.tmpl.html", templateData)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippet = snippet
	app.render(w, r, http.StatusOK, "view.tmpl.html", templateData)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)

	templateData.Form = snippetCreateForm{
		ExpiresAt: 365,
	}
	app.render(w, r, http.StatusOK, "create.tmpl.html", templateData)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, r, err)
	}
	expiresAfterDays := r.PostForm.Get("expires")
	expiresAfterDaysInt, err := strconv.Atoi(expiresAfterDays)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		ExpiresAt:   expiresAfterDaysInt,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if expiresAfterDaysInt != 1 && expiresAfterDaysInt != 7 && expiresAfterDaysInt != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, expiresAfterDaysInt)
	if err != nil {
		app.serverError(w, r, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
