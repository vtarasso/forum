package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	post_id, err := strconv.Atoi(url[len(url)-1])
	if err != nil || len(url) != 3 || post_id < 1 {
		app.notFound(w, r)
		return
	}
	status := r.Context().Value("token").(*Data)
	user, _ := app.users.GetUser(status.Token)
	if r.Method != http.MethodPost {
		s, err := app.snippets.GetById(post_id)
		if err != nil {
			if errors.Is(err, ErrNoRecord) {
				app.notFound(w, r)
			} else {
				app.serverError(w, err, r)
			}
			return
		}
		snippetscom, err := app.snippets.GetComments(post_id)
		if err != nil {
			if errors.Is(err, ErrNoRecord) {
				app.notFound(w, r)
			} else {
				app.serverError(w, err, r)
			}
			return
		}
		data := app.newTemplateData(r)
		data.SnippetsCom = snippetscom
		data.UserData = user
		data.Snippet = s
		app.render(w, http.StatusOK, "showpage.html", data, r)
	} else if r.Method == http.MethodPost {
		if user == nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		err := r.ParseForm()
		if err != nil {
			app.serverError(w, err, r)
			return
		}
		s, err := app.snippets.GetById(post_id)
		if err != nil {
			app.serverError(w, err, r)
			return
		}

		formComment := FormComment{
			UserId:    status.UserID,
			UserName:  user.Name,
			PostID:    post_id,
			Content:   r.PostForm.Get("content"),
			Validator: Validator{},
		}
		data := app.newTemplateData(r)
		formComment.Validator.CheckField(NotBlank(formComment.Content), "content", ErrBlank)
		if !formComment.Validator.Valid() {
			data.Snippet = s
			data.UserData = user
			data.SnippetsCom, err = app.snippets.GetComments(post_id)
			if err != nil {
				app.serverError(w, err, r)
				return
			}
			data.FormCom = formComment
			app.render(w, http.StatusUnprocessableEntity, "showpage.html", data, r)
			return
		}
		err = app.snippets.InsertComments(&formComment)
		if err != nil {
			app.serverError(w, err, r)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", post_id), http.StatusSeeOther)
	}
}

func (app *application) myPosts(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value("token").(*Data)

	if r.Method != http.MethodPost {
		dataposts := app.newTemplateData(r)
		datacat, err := app.snippets.GetAllCatigories(dataposts)
		if err != nil {
			app.serverError(w, err, r)
		}
		user, err := app.users.GetUser(ctxData.Token)
		if err != nil {
			app.serverError(w, err, r)
			return
		}

		s, err := app.snippets.ReturMyPosts(user.Name)
		if err != nil {
			if errors.Is(err, ErrNoRecord) {
				app.notFound(w, r)
			} else {
				app.serverError(w, err, r)
			}
			return
		}

		dataposts.Snippets = s
		dataposts.UserData = user
		dataposts.Catigories = datacat

		var snip []*Snippet

		for w := len(dataposts.Snippets) - 1; w >= 0; w-- {
			snip = append(snip, dataposts.Snippets[w])
		}
		dataposts.Snippets = snip
		
		app.render(w, http.StatusOK, "myposts.html", dataposts, r)

	}
}
