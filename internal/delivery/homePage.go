package delivery

import (
	"log"
	"net/http"
	"text/template"

	"forum/internal/models"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(key).(models.User)
	if !ok {
		ts, err := template.ParseFiles("./ui/templates/index.html")
		if err = ts.Execute(w, nil); err != nil {
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		return
	}
	if r.URL.Path != "/" {
		log.Println("home page: wrong url")
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		log.Println("home page:Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	ts, err := template.ParseFiles("./ui/templates/index.html")
	if err != nil {
		log.Println(err.Error())
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if err = ts.Execute(w, user); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}
