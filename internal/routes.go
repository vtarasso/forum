package internal

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.requireAuthentication(app.createSnippetPost))
	mux.HandleFunc("/user/signup", app.AuthCheck(app.userSignupPost))
	mux.HandleFunc("/user/login", app.AuthCheck(app.userLoginPost))
	mux.HandleFunc("/user/logout", app.requireAuthentication(app.userLogout))
	mux.HandleFunc("/likePost", app.requireAuthentication(app.likePost))
	mux.HandleFunc("/dislikePost", app.requireAuthentication(app.dislikePost))
	mux.HandleFunc("/likeComment", app.requireAuthentication(app.likeComment))
	mux.HandleFunc("/dislikeComment", app.requireAuthentication(app.dislikeComment))
	mux.HandleFunc("/snippet/myposts", app.requireAuthentication(app.myPosts))
	mux.HandleFunc("/snippet/liked", app.requireAuthentication(app.likedPosts))

	fileServer := http.FileServer(http.Dir("assets/static"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return app.recoverPanic((app.logRequest(secureHeaders(app.myMiddleware(mux)))))
}
