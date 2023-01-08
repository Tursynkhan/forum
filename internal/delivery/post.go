package delivery

import (
	"log"
	"net/http"
	"text/template"

	"forum/internal/models"
)

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create-post" {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if r.Method == "GET" {
		ts, err := template.ParseFiles("./ui/templates/createPost.html")
		if err != nil {
			log.Printf("Create Post: Execute:%v", err)
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
		err := r.ParseForm()
		if err != nil {
			log.Println("error parse form :", err)
			return
		}
		title := r.PostFormValue("title")
		content := r.PostFormValue("content")

		newPost := models.Post{
			Title:   title,
			Content: content,
		}

		if err = h.services.Post.CreatePost(newPost); err != nil {
			log.Printf("Post: Create Post: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		log.Println("Create Post: Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
}
