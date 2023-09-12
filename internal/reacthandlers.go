package internal

import (
	"net/http"
	"strconv"
)

func (app *application) likedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		status := r.Context().Value("token").(*Data)
		data := app.newTemplateData(r)
		var err error
		if err != nil {
			app.serverError(w, err, r)
			return
		}
		data.Snippets, err = app.snippets.GetLikedPost(status.UserID)
		if err != nil {
			app.serverError(w, err, r)
			return
		}
		data.UserData, err = app.users.GetUser(status.Token)
		if err != nil {
			app.serverError(w, err, r)
			return
		}

		var snip []*Snippet

		for w := len(data.Snippets) - 1; w >= 0; w-- {
			snip = append(snip, data.Snippets[w])
		}
		data.Snippets = snip
		app.render(w, http.StatusOK, "likedposts.html", data, r)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *application) dislikePost(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("token").(*Data)
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusUnprocessableEntity)
	}
	err = app.reactions.DislikePost(data.UserID, postID)
	if err != nil {
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusUnprocessableEntity)
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *application) likePost(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("token").(*Data)
	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusUnprocessableEntity)
	}
	err = app.reactions.LikePost(data.UserID, postID)
	if err != nil {
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusUnprocessableEntity)
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *application) likeComment(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("token").(*Data)
	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || commentID < 1 {
		app.notFound(w, r)
		return
	}
	_ = app.reactions.LikeComment(status.UserID, commentID)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *application) dislikeComment(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("token").(*Data)
	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || commentID < 1 {
		app.notFound(w, r)
		return
	}
	_ = app.reactions.DislikeComment(status.UserID, commentID)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
