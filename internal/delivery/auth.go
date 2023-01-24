package delivery

import (
	"errors"
	"fmt"
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
	// user, ok := h.ctx.Context().Value(key).(models.User)
	// if !ok {
	// 	h.errorHandler(w, http.StatusBadRequest, "You already in")
	// 	return
	// }
	if r.URL.Path != "/auth/signup" {
		log.Println("Sign Up:Wrong URL Path")
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method == "GET" {
		Form := NewForm{
			Form: *forms.New(nil),
		}
		if err := h.tmpl.ExecuteTemplate(w, "signUp.html", Form); err != nil {
			log.Println(err.Error())
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		name := r.FormValue("name")
		email := r.FormValue("email")
		psw := r.FormValue("psw")
		pswRepeat := r.FormValue("psw-repeat")

		newUser := models.User{
			Username:       name,
			Email:          email,
			Password:       psw,
			RepeatPassword: pswRepeat,
		}
		if err := h.services.Autorization.CreateUser(newUser); err != nil {
			form := forms.New(r.PostForm)
			form.Required("name", "email", "psw")
			log.Printf("Sign Up: Create User: %v", err)
			if errors.Is(err, service.ErrInvalidEmail) {
				form.ErrorField("email")
			} else if errors.Is(err, service.ErrInvalidPassword) {
				form.ErrorMatchesPattern("psw")
			} else if errors.Is(err, service.ErrInvalidUsername) {
				form.ErrorField("name")
			} else if errors.Is(err, service.ErrUserExist) {
				form.IsExist("name", "User exist")
			} else if errors.Is(err, service.ErrUserNotFound) {
				form.IsExist("name", "User Not Found")
			} else if errors.Is(err, service.ErrPasswdNotMatch) {
				form.IsExist("psw", "Password doesn't match")
			} else if errors.Is(err, service.ErrEmailExist) {
				form.IsExist("email", "Email exist")
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
	} else {
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
	if r.Method == "GET" {
		Form := NewForm{
			Form: *forms.New(nil),
		}
		if err := h.tmpl.ExecuteTemplate(w, "signIn.html", Form); err != nil {
			log.Println(err.Error())
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		name := r.FormValue("name")
		psw := r.FormValue("psw")

		form := forms.New(r.PostForm)
		sessionToken, expireTime, err := h.services.GenerateToken(name, psw)
		if err != nil {
			log.Printf("Sign In: Generate Token:%v", err)
			if errors.Is(err, service.ErrUserNotFound) {
				form.Errors.Add("generic", "Username or Password is incorrect")
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
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expireTime,
			Path:    "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/logout" {
		log.Println("Sign In : Wrong URL Path")
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method == "POST" {
		var err error
		token, err := r.Cookie("session_token")
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
			Name:    "session_token",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		}
		http.SetCookie(w, c)

		fmt.Println("after setcookie logout")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.Println("Logout : Method Not Allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
