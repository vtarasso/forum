package internal

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Определяем новый тип Validator, который содержит карту ошибок проверки для нашего
// поля формы.
// Добавьте новое поле NonFieldErrors []string в структуру, которое мы будем использовать
// для хранения всех ошибок валидации, которые не связаны с конкретным полем формы.
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid() возвращает true, если карта FieldErrors не содержит записей.
// Обновите метод Valid(), чтобы он также проверял, что срез NonFieldErrors пуст.

func (v *Validator) Valid() bool {
	if len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0 {
		return true
	}
	return false
}

// AddFieldError() добавляет сообщение об ошибке в карту FieldErrors
// (если для данного ключа еще не существует записи).
func (v *Validator) AddFieldError(key, message string) {
	// Примечание. Сначала нам нужно инициализировать карту, если она еще не инициализирована.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Создайте вспомогательную функцию AddNonFieldError() для добавления сообщений об ошибке в новый срез NonFieldErrors.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField() добавляет сообщение об ошибке в карту FieldErrors, только если
// проверка проверки не "ok".
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() возвращает true, если значение не является пустой строкой.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() возвращает true, если значение содержит не более n символов.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func CatIsNill(catId []int) bool {
	return catId != nil
}

func isValidEmail(email string) bool {
	rxEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidName(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < 5 || length > 15 {
		return false
	}
	usernameConvention := "^[A-Za-z][A-Za-z0-9_]{4,14}$"
	re, _ := regexp.Compile(usernameConvention)
	return re.MatchString(name)
}

func isValidPassword(pass string) bool {
	tests := []string{".{8,20}", "[A-Z]", "[a-z]", "[0-9]", "[!,@,#,$,%,^,&,*,(,),_,-,+,=,?,|,/,;,:,{,},.,,]"}
	for _, test := range tests {
		valid, _ := regexp.MatchString(test, pass)
		if !valid {
			return false
		}
	}
	return true
}
