package delivery

import (
	"forum/internal/forms"
	"forum/internal/models"
	"log"
	"net/http"
	"text/template"
	"time"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		log.Println("Sign Up:Wrong URL Path")
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method == "GET" {
		ts, err := template.ParseFiles("./ui/templates/signUp.html")
		if err != nil {
			log.Printf("Sign Up: Execute:%v", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if r.Method == "POST" {

		if err := r.ParseForm(); err != nil {
			h.errorHandler(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		name := r.FormValue("name")
		email := r.FormValue("email")
		psw := r.FormValue("psw")
		pswRepeat := r.FormValue("psw-repeat")

		form := forms.New(r.PostForm)
		form.Required("name", "email", "password")
		form.MatchesPattern("email", forms.EmailRX)
		form.Minlength("psw", 10)

		if !form.Valid() {
			ts, err := template.ParseFiles("./ui/templates/signUp.html")
			if err != nil {
				log.Printf("Sign Up: Execute:%v", err)
				h.errorHandler(w, http.StatusInternalServerError, err.Error())
				return
			}
			type NewForm struct {
				Form *forms.Form
			}
			Form := NewForm{
				Form: form,
			}
			err = ts.Execute(w, Form)
			if err != nil {
				log.Println(err.Error())
				h.errorHandler(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		newUser := models.User{
			Username:       name,
			Email:          email,
			Password:       psw,
			RepeatPassword: pswRepeat,
		}
		if err := h.services.Autorization.CreateUser(newUser); err != nil {
			log.Printf("Sign Up: Create User: %v", err)
			h.errorHandler(w, http.StatusForbidden, err.Error())
			return
		}
		http.Redirect(w, r, "/auth/signin", 301)
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
		ts, err := template.ParseFiles("./ui/templates/signIn.html")
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err = ts.Execute(w, nil); err != nil {
			log.Printf("Sign In: Execute:%v", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		psw := r.FormValue("psw")
		sessionToken, expiresTime, err := h.services.GenerateToken(name, psw)
		if err != nil {
			log.Printf("Sign In: Generate Token:%v", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiresTime,
			Path:    "/",
		})
		http.Redirect(w, r, "/", 301)
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/auth/logout" {
	// 	log.Println("Sign In : Wrong URL Path")
	// 	h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	// 	return
	// }
	if r.Method == "GET" {
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
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		http.Redirect(w, r, "/", 301)
	} else {
		log.Println("Logout : Method Not Allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
