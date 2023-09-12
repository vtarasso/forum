package internal

import (
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("token").(*Data)
	if r.URL.Path != "/" {
		app.notFound(w, r)
		return
	}
	data := app.newTemplateData(r)
	datacat, err := app.snippets.GetAllCatigories(data)
	user, _ := app.users.UserbyToken(status.Token)
	if err != nil {
		app.serverError(w, err, r)
	}
	switch r.Method {
	case http.MethodGet:
		s, err := app.snippets.Latest()
		if err != nil {
			app.serverError(w, err, r)
			return
		}

		data.UserData = user
		data.Snippets = s
		data.Catigories = datacat

		var snip []*Snippet

		for w := len(data.Snippets) - 1; w >= 0; w-- {
			snip = append(snip, data.Snippets[w])
		}
		data.Snippets = snip

		app.render(w, http.StatusOK, "homepage.html", data, r)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			app.errorHandler(w, r, http.StatusBadRequest)
			return
		}
		var catID []int
		for _, value := range r.PostForm["cat"] {
			number, err := strconv.Atoi(value)
			if err != nil {
				app.serverError(w, err, r)
				return
			}
			catID = append(catID, number)
		}
		data.Snippets, err = app.snippets.GetFilterPost(catID, status.UserID)
		if err != nil {
			app.serverError(w, err, r)
			return
		}
		if catID != nil {
			data.Snippets, err = app.snippets.GetFilterPost(catID, status.UserID)
			if err != nil {
				app.serverError(w, err, r)
				return
			}
		} else {
			data.Snippets, err = app.snippets.Latest()
			if err != nil {
				app.serverError(w, err, r)
				return
			}
		}
		data.UserData = user

		app.render(w, http.StatusOK, "homepage.html", data, r)
	}
}
