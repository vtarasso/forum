package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

// type contextKey string

// const ctxKey contextKey = "data"

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, fmt.Errorf("%s", w), r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) myMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		data := &Data{}
		switch err {
		case http.ErrNoCookie:
			data.UserID = 0
		case nil:
			user, err := app.users.UserbyToken(cookie.Value)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				app.serverError(w, err, r)
				return
			}
			if user != nil {
				data.IsAuthenticated = true
				data.Token = cookie.Value
				data.UserID = user.ID
			}
		default:
			app.serverError(w, err, r)
			return
		}
		ctx := context.WithValue(r.Context(), "token", data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value("token").(*Data)
		if !data.IsAuthenticated {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	}
}

func (app *application) AuthCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		next.ServeHTTP(w, r)
	}
}
