package delivery

import (
	"forum/internal/models"
	"log"
	"net/http"
	"text/template"
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
		categories, err := h.services.GetAllCategories()
		if err != nil {
			log.Println("home page : get all categories", err)
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
			posts, err = h.services.GetPostByFilter(r.URL.Query())
			if err != nil {
				log.Println("home page : GetPostByFilter : ", err)
				h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		info := models.Info{
			Posts:    posts,
			Category: categories,
		}
		ts, err := template.ParseFiles("./ui/templates/index.html")
		if err = ts.Execute(w, info); err != nil {
			log.Printf("homepage : execute : %v", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		return
	}

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
		posts, err = h.services.GetPostByFilter(r.URL.Query())
		if err != nil {
			log.Println("home page : GetPostByFilter : ", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	info := models.Info{
		Posts:    posts,
		User:     user,
		Category: categories,
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
