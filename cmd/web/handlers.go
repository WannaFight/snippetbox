package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/WannaFight/snippetbox/internal/models"
	"github.com/WannaFight/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

const (
	flashSessionKey      = "flash"
	authUserIDSessionKey = "authenticatedUserID"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	if err := app.decodePostForm(r, &form); err != nil {
		app.badRequest(w)
		return
	}

	titleMaxLen := 100
	expiresValues := []int{1, 7, 365}

	form.CheckField(validator.NotBlank(form.Title), "title", validator.BlankStringValidationError)
	form.CheckField(validator.MaxChars(form.Title, titleMaxLen), "title", fmt.Sprintf(validator.TextTooLongValidationError, titleMaxLen))
	form.CheckField(validator.NotBlank(form.Content), "content", validator.BlankStringValidationError)
	form.CheckField(validator.PermittedValue(form.Expires, expiresValues...), "expires", fmt.Sprintf(validator.ChoiceValidationError, expiresValues))

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), flashSessionKey, "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	form := new(userSignupForm)

	r.ParseForm()
	fmt.Println("post form\n", r.PostForm)

	if err := app.decodePostForm(r, &form); err != nil {
		app.badRequest(w)
		return
	}

	passwordLen := 8

	fmt.Println("form\n", form)

	form.CheckField(validator.NotBlank(form.Name), "name", validator.BlankStringValidationError)
	form.CheckField(validator.NotBlank(form.Email), "email", validator.BlankStringValidationError)
	form.CheckField(validator.ValidEmail(form.Email), "email", validator.NotValidEmailValidationError)
	form.CheckField(validator.NotBlank(form.Password), "password", validator.BlankStringValidationError)
	form.CheckField(validator.MinChars(form.Password, passwordLen), "password", fmt.Sprintf(validator.TextTooShortValidationError, passwordLen))

	fmt.Println("form\n", form)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	if err := app.users.Insert(form.Name, form.Email, form.Password); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", validator.DuplicateEmailValidationError)

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmplva", data)
			return
		}
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), flashSessionKey, "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	form := new(userLoginForm)

	if err := app.decodePostForm(r, &form); err != nil {
		app.badRequest(w)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", validator.BlankStringValidationError)
	form.CheckField(validator.ValidEmail(form.Email), "email", validator.NotValidEmailValidationError)
	form.CheckField(validator.NotBlank(form.Password), "password", validator.BlankStringValidationError)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
			return
		}
		app.serverError(w, err)
		return
	}

	if err = app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), authUserIDSessionKey, id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), authUserIDSessionKey)
	app.sessionManager.Put(r.Context(), flashSessionKey, "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
