package internal

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// Добавляем статус аутентификации в данные шаблона.
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// Обработка ошибки
		return false
	}
	sessionID := cookie.Value
	_, err = app.users.UserbyToken(sessionID)
	if err != nil {
		return false
	}
	return true
}

func (app *application) serverError(w http.ResponseWriter, err error, r *http.Request) {
	data := app.newTemplateData(r)

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	dataerr := &Data{
		IntErr: http.StatusInternalServerError,
	}
	data.Data = dataerr

	app.errorHandler(w, r, data.Data.IntErr)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	dataerr := &Data{
		IntErr: http.StatusNotFound,
	}
	data.Data = dataerr
	app.errorHandler(w, r, data.Data.IntErr)
}

func (app *application) render(w http.ResponseWriter, status int, name string, td *templateData, r *http.Request) {
	ts, ok := app.templateCache[name]
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err, r)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) renderErr(w http.ResponseWriter, status int, name string, td *templateData, r *http.Request) error {
	ts, ok := app.templateCache[name]
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return ErrInvalidObjectId
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err, r)
		return err
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
	return err
}
