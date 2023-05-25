package delivery

import (
	"errors"
	"forum/internal/forms"
	"forum/internal/models"
	"forum/internal/service"
	"log"
	"net/http"
	"time"
)

type NewForm struct {
	Form forms.Form
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user != (models.User{}) {
		h.errorHandler(w, http.StatusBadRequest, "you already in")
		return
	}
	if r.URL.Path != "/auth/signup" {
		log.Println("Sign Up:Wrong URL Path")
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {

	case http.MethodGet:
		Form := NewForm{
			Form: *forms.New(nil),
		}
		if err := h.tmpl.ExecuteTemplate(w, "signUp.html", Form); err != nil {
			log.Println(err.Error())
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:

		err := r.ParseForm()
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		name, ok := r.Form["name"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "email field not found")
			return
		}
		email, ok := r.Form["email"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "username field not found")
			return
		}
		psw := r.Form["psw"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "password field not found")
			return
		}
		pswRepeat := r.Form["psw-repeat"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "psw-repeat field not found")
			return
		}

		newUser := models.User{
			Username:       name[0],
			Email:          email[0],
			Password:       psw[0],
			RepeatPassword: pswRepeat[0],
			RoleID:         1,
		}
		if err := h.services.Autorization.CreateUser(newUser); err != nil {
			form := forms.New(r.PostForm)
			form.Required("name", "email", "psw")
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Sign Up: Create User: %v", err)
			if errors.Is(err, service.ErrInvalidEmail) {
				form.ErrorField("email")
				w.WriteHeader(http.StatusBadRequest)
			} else if errors.Is(err, service.ErrInvalidPassword) {
				form.ErrorMatchesPattern("psw")
				w.WriteHeader(http.StatusBadRequest)

			} else if errors.Is(err, service.ErrInvalidUsername) {
				form.ErrorField("name")
				w.WriteHeader(http.StatusBadRequest)

			} else if errors.Is(err, service.ErrUserExist) {
				form.IsExist("name", "User exist")
				w.WriteHeader(http.StatusBadRequest)

			} else if errors.Is(err, service.ErrUserNotFound) {
				form.IsExist("name", "User Not Found")
				w.WriteHeader(http.StatusBadRequest)

			} else if errors.Is(err, service.ErrPasswdNotMatch) {
				form.IsExist("psw", "Password doesn't match")
				w.WriteHeader(http.StatusBadRequest)

			} else if errors.Is(err, service.ErrEmailExist) {
				form.IsExist("email", "Email exist")
				w.WriteHeader(http.StatusBadRequest)

			} else {
				h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			if !form.Valid() {
				Form := NewForm{
					Form: *form,
				}
				if err := h.tmpl.ExecuteTemplate(w, "signUp.html", Form); err != nil {
					log.Println(err.Error())
					h.errorHandler(w, http.StatusInternalServerError, err.Error())
					return
				}
				return
			}
		}
		http.Redirect(w, r, "/auth/signin", http.StatusSeeOther)
	default:
		log.Println("Sign Up: Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return

	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signin" {
		log.Println("Sign In:Wrong URL Path")
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {

	case http.MethodGet:
		Form := NewForm{
			Form: *forms.New(nil),
		}
		if err := h.tmpl.ExecuteTemplate(w, "signIn.html", Form); err != nil {
			log.Println(err.Error())
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		name, ok := r.Form["name"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "name field not found")
			return
		}
		psw, ok := r.Form["psw"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "psw field not found")
			return
		}
		form := forms.New(r.PostForm)
		sessionToken, expireTime, err := h.services.GenerateToken(name[0], psw[0])
		if err != nil {
			log.Printf("Sign In: Generate Token:%v", err)
			if errors.Is(err, service.ErrUserNotFound) {
				w.WriteHeader(http.StatusBadRequest)
				form.Errors.Add("generic", "Username doesn't exist or Password is incorrect")
				Form := NewForm{
					Form: *form,
				}
				if err := h.tmpl.ExecuteTemplate(w, "signIn.html", Form); err != nil {
					log.Println(err.Error())
					h.errorHandler(w, http.StatusInternalServerError, err.Error())
					return
				}
				return
			}
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   sessionToken,
			Expires: expireTime,
			Path:    "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		log.Println("Sign Up: Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusBadRequest, "can't log-out,without log-in")
		return
	}
	if r.URL.Path != "/auth/logout" {
		log.Println("Sign In : Wrong URL Path")
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {
	case http.MethodGet:
		var err error
		token, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
		}
		err = h.services.DeleteToken(token.Value)
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		c := &http.Cookie{
			Name:    "token",
			Value:   "",
			Path:    "/",
			Expires: time.Now(),
		}
		http.SetCookie(w, c)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		log.Println("Logout : Method Not Allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
