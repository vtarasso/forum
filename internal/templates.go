package internal

import (
	"path/filepath"
	"text/template"
	"time"
)

// Создаем тип templateData, который будет действовать как хранилище для
// любых динамических данных, которые нужно передать HTML-шаблонам.
// На данный момент он содержит только одно поле, но мы добавим в него другие
// по мере развития нашего приложения.

// Добавляем поле Snippets в структуру templateData
// Добавляем поле CurrentYear в структуру templateData.

type templateData struct {
	CurrentYear     int
	UserData        *User
	Snippet         *Snippet
	Snippets        []*Snippet
	SnippetCom      *SnippetCom
	SnippetsCom     []*SnippetCom
	FormCom         FormComment
	FormUsers       userSignupForm
	FormLogin       userLoginForm
	Catigories      []string
	Data            *Data
	IsAuthenticated bool
	ErrorStruct     ErrorStruct
}

// Создаем функцию humanDate, которая возвращает красиво отформатированную строку
// представление объекта time.Time.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Инициализируем новую карту, которая будет хранить кэш.
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, err
	}
	// Перебираем файл шаблона от каждой страницы.
	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
