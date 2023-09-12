package internal

import (
	"errors"
	"net/http"
)

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		data := app.newTemplateData(r)
		data.FormUsers = userSignupForm{}
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data, r)
	} else if r.Method != http.MethodGet {
		// Declare an zero-valued instance of our userSignupForm struct.
		// var formUsers userSignupForm

		// Парсим данные формы в структуру userSignupForm.
		err := r.ParseForm()
		if err != nil {
			app.errorHandler(w, r, http.StatusBadRequest)
			return
		}
		// Объявляем экземпляр структуры userSignupForm со значениями по умолчанию.
		formUsers := &userSignupForm{
			Name:      r.PostForm.Get("name"),
			Email:     r.PostForm.Get("email"),
			Password:  r.PostForm.Get("password"),
			Validator: Validator{},
		}

		// Проверяем содержимое формы с помощью наших вспомогательных функций.
		formUsers.Validator.CheckField(NotBlank(formUsers.Name), "name", ErrBlank)
		formUsers.Validator.CheckField(isValidName(formUsers.Name), "name", ErrCorrectName)
		formUsers.Validator.CheckField(NotBlank(formUsers.Email), "email", ErrBlank)
		formUsers.Validator.CheckField(isValidEmail(formUsers.Email), "email", ErrEmail)
		formUsers.Validator.CheckField(NotBlank(formUsers.Password), "password", ErrBlank)
		formUsers.Validator.CheckField(isValidPassword(formUsers.Password), "password", ErrPass)

		// Если есть ошибки, снова отображаем форму регистрации вместе с кодом состояния 422.
		if !formUsers.Validator.Valid() {
			data := app.newTemplateData(r)
			data.FormUsers = *formUsers
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data, r)
			return
		}

		err = app.users.SignUp(formUsers)
		if err != nil {
			switch err {
			case ErrDuplicateEmail:
				formUsers.Validator.AddFieldError("email", ErrUsedEmail)
				user, err := app.users.UserByName(formUsers.Name)
				if err != nil && !errors.Is(err, ErrNoRecord) {
					return
				}
				if user.Name == formUsers.Name {
					formUsers.Validator.AddFieldError("name", ErrUsedName)
				}
			case ErrDuplicateName:
				formUsers.Validator.AddFieldError("name", ErrUsedName)
			default:
				app.serverError(w, err, r)
			}
			data := app.newTemplateData(r)
			data.FormUsers = *formUsers
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data, r)
		}
		// Затем перенаправляем пользователя на страницу входа.
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

// Create a new userLoginForm struct.

// Update the handler so it displays the login page.
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		data := app.newTemplateData(r)
		data.FormLogin = userLoginForm{}
		app.render(w, http.StatusOK, "login.html", data, r)
	} else if r.Method != http.MethodGet {
		// Парсим данные формы в структуру userLoginForm.
		err := r.ParseForm()
		if err != nil {
			app.errorHandler(w, r, http.StatusBadRequest)
			return
		}
		formLogin := userLoginForm{
			Email:     r.PostForm.Get("email"),
			Password:  r.PostForm.Get("password"),
			Validator: Validator{},
		}

		formLogin.Validator.CheckField(NotBlank(formLogin.Email), "email", ErrBlank)
		formLogin.Validator.CheckField(isValidEmail(formLogin.Email), "email", ErrValidEmail)
		formLogin.Validator.CheckField(NotBlank(formLogin.Password), "password", ErrBlank)
		if !formLogin.Validator.Valid() {
			data := app.newTemplateData(r)
			data.FormLogin = formLogin
			app.render(w, http.StatusUnprocessableEntity, "login.html", data, r)
			return
		}

		id, err := app.users.Authenticate(formLogin.Email, formLogin.Password)
		datauser := Data{
			UserID: id,
		}
		datauser.UserID = id
		if err != nil {
			if errors.Is(err, ErrInvalidCredentials) {
				formLogin.Validator.AddNonFieldError("Email or password is incorrect")
				data := app.newTemplateData(r)
				data.FormLogin = formLogin
				app.render(w, http.StatusUnprocessableEntity, "login.html", data, r)
			} else {
				app.serverError(w, err, r)
			}
			return
		}
		token, err := GenerateToken()
		if err != nil {
			app.serverError(w, err, r)
			return
		}
		err = app.users.SetToken(id, *token)
		if err != nil {
			app.serverError(w, err, r)
			return
		}
		cookie := &http.Cookie{
			Path:  "/",
			Name:  "session_token",
			Value: *token,
		}
		http.SetCookie(w, cookie)

		// Перенаправляем пользователя на страницу создания фрагмента кода.
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

// Для удаления tokena пользователя из BD и закрытия сессий
func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	ctxData := r.Context().Value("token").(*Data)
	app.users.RemoveToken(ctxData.Token)
	cookie := &http.Cookie{
		Path: "/",
		Name: "session_token",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
