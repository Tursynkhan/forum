package delivery

import (
	"forum/internal/models"
	"log"
	"net/http"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	var posts []models.PostInfo
	var err error
	if r.URL.Path != "/" {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		log.Println("home page : Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := r.Context().Value(key).(models.User)
	categories, err := h.services.GetAllCategories()
	if err != nil {
		log.Println("home page : get all categories :", err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if len(r.URL.Query()) == 0 {
		posts, err = h.services.GetAllPosts()
		if err != nil {
			log.Println("home page : get all posts : ", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	} else {
		posts, err = h.services.GetPostByFilter(r.URL.Query(), user)
		if err != nil {
			log.Println("home page : GetPostByFilter : ", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	notifications, err := h.services.GetAllNotification(user)
	if err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	info := models.Info{
		Posts:         posts,
		User:          user,
		Category:      categories,
		Notifications: notifications,
	}
	if err := h.tmpl.ExecuteTemplate(w, "index.html", info); err != nil {
		log.Println(err.Error())
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
}
