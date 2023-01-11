package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		// categories := r.Form["categories"]
		newPost := models.Post{
			Title:   title,
			Content: content,
		}

		if err = h.services.Post.CreatePost(newPost); err != nil {
			log.Printf("Post: Create Post: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		fmt.Println("Posts:", newPost)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		log.Println("Create Post: Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		user, ok := r.Context().Value(key).(models.User)
		if !ok {
		}
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		fmt.Println("Get Post: This is Id:", id)
		post, err := h.services.Post.GetPost(id)
		fmt.Println("Get Post: This is post:", post)
		if err != nil {
			log.Printf("Post: getPost: %v", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		ts, err := template.ParseFiles("./ui/templates/post.html")
		if err != nil {
			log.Printf("Get Post: Execute:%v", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		info := models.Info{
			User: user,
			Post: post,
		}
		err = ts.Execute(w, info)
		if err != nil {
			log.Println(err.Error())
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		log.Println("Get Post: Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}
