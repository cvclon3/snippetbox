package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"errors"

	"snippetbox.cvclon3.net/internal/models"
)


func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%+v\n", snippet)
	// }

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// !!! tmpl = ts in the book
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		// app.errorLog.Println(err.Error())
		app.serverError(w, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
	}


	data := &templateData{
		Snippets: snippets,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
	}

	// w.Write([]byte("Hello from App"))
}


func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippet: snippet,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
	
	// fmt.Fprintf(w, "View snippet for ID=%d", id)
	// fmt.Fprintf(w, "%+v", snippet)
}


func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // POST
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405)
		// w.Write([]byte("Method now allowed"))
		// http.Error(w, "Method now allowed", http.StatusMethodNotAllowed) //405
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 shall"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

	// w.Write([]byte("Create snippet"))
}
