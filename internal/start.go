package internal

import (
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	"forum/internal/sqlitedb"
)

type application struct {
	reactions     *ReactionsModel
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *SnippetModel
	users         *UserModel
	templateCache map[string]*template.Template
}

func StartGo() {
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := sqlitedb.CreateDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./assets/templates/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		reactions:     &ReactionsModel{DB: db},
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &SnippetModel{DB: db},
		users:         &UserModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("listening on http://localhost" + *addr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
