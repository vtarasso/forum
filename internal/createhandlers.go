package internal

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	status := r.Context().Value("token").(*Data)
	data := app.newTemplateData(r)
	datacat, err := app.snippets.GetAllCatigories(data)
	if err != nil {
		app.serverError(w, err, r)
		return
	}
	if r.Method != http.MethodPost {
		data.Catigories = datacat
		user, _ := app.users.UserbyToken(status.Token)
		post := Snippet{
			UserName: user.Name,
		}
		data.Snippet = &post
		app.render(w, http.StatusUnprocessableEntity, "createpage.html", data, r)
	} else if r.Method == http.MethodPost {
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
		user, _ := app.users.UserbyToken(status.Token)
		post := Snippet{
			UserName:  user.Name,
			Title:     r.PostForm.Get("title"),
			Content:   r.PostForm.Get("content"),
			Validator: Validator{},
			CatID:     catID,
		}
		post.Validator.CheckField(NotBlank(post.Title), "title", ErrBlank)
		post.Validator.CheckField(MaxChars(post.Title, 100), "title", ErrMaxChars)
		post.Validator.CheckField(NotBlank(post.Content), "content", ErrBlank)
		post.Validator.CheckField(CatIsNill(post.CatID), "cat", ErrChoiceCategory)

		if !post.Validator.Valid() {
			data.Catigories = datacat
			data.UserData = user
			data.Snippet = &post
			data.IsAuthenticated = status.IsAuthenticated
			app.render(w, http.StatusUnprocessableEntity, "createpage.html", &templateData{
				Snippet:         &post,
				UserData:        user,
				IsAuthenticated: status.IsAuthenticated,
				Catigories:      datacat,
			}, r)
			return
		}
		id, err := app.snippets.Create(&post)
		if err != nil {
			app.serverError(w, err, r)
			return
		}

		// Перенаправляем пользователя на соответствующую страницу заметки.
		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	}
}

func (app *application) errorHandler(w http.ResponseWriter, r *http.Request, code int) {
	data := app.newTemplateData(r)
	// data.ErrorStruct = ErrorStruct{}
	Res := &ErrorStruct{
		Status: code,
		Text:   http.StatusText(code),
	}
	data.ErrorStruct = *Res
	err := app.renderErr(w, http.StatusUnprocessableEntity, "error.html", data, r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprintln(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}
