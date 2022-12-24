package delivery

import (
	"log"
	"net/http"
	"text/template"

	"forum/internal/models"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/signup" {
		h.ErrorHandler(w, r, errStatus{http.StatusNotFound, http.StatusText(http.StatusNotFound)})
	}
	if r.Method == "GET" {
		ts, err := template.ParseFiles("./ui/templates/signUp.html")
		if err != nil {
			h.ErrorHandler(w, r, errStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			h.ErrorHandler(w, r, errStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		}
	} else if r.Method == "POST" {

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
			h.ErrorHandler(w, r, errStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
			return
		}
		http.Redirect(w, r, "/auth/signin", 303)
	}
}

type signInInput struct {
	Username string
	Password string
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ts, err := template.ParseFiles("./ui/templates/signIn.html")
		if err != nil {
			h.ErrorHandler(w, r, errStatus{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		}
		err = ts.Execute(w, nil)
	} else if r.Method == "POST" {
	}
}
