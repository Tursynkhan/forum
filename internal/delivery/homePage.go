package delivery

import (
	"log"
	"net/http"
	"text/template"

	"forum/internal/models"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	var posts []models.PostInfo
	var err error
	if r.URL.Path != "/" {
		log.Println("home page : wrong url")
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		log.Println("home page : Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user, ok := r.Context().Value(key).(models.User)
	if !ok {
		posts, err = h.services.GetAllPosts()
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		info := models.Info{
			Posts: posts,
		}
		ts, err := template.ParseFiles("./ui/templates/index.html")
		if err = ts.Execute(w, info); err != nil {
			log.Printf("homepage : execute : %v", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		return
	}
	posts, err = h.services.GetAllPosts()
	if err != nil {
		log.Println("home page : get all posts", err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	info := models.Info{
		Posts: posts,
		User:  user,
	}
	ts, err := template.ParseFiles("./ui/templates/index.html")
	if err != nil {
		log.Println(err.Error())
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if err = ts.Execute(w, info); err != nil {
		log.Println("home page : execute error :", err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}
