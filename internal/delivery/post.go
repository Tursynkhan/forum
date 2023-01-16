package delivery

import (
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
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
		user, ok := r.Context().Value(key).(models.User)
		if !ok {
			h.errorHandler(w, http.StatusInternalServerError, "Unauthorized")
			return
		}
		err := r.ParseForm()
		if err != nil {
			log.Println("error parse form :", err)
			return
		}
		title := r.PostFormValue("title")
		content := r.PostFormValue("content")
		categories := r.Form["categories"]
		newPost := models.Post{
			UserID:  user.ID,
			Title:   title,
			Content: content,
		}
		postId, err := h.services.Post.CreatePost(newPost)
		if err != nil {
			log.Printf("Post: Create Post: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if err = h.services.Post.CreatePostCategory(postId, categories); err != nil {
			log.Printf("Post: Create PostCategory : %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
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
			h.errorHandler(w, http.StatusInternalServerError, "Unauthorized")
			return
		}

		id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/get-post/"))
		fmt.Println("Get Post: This is Id:", id)

		post, err := h.services.Post.GetPost(id)
		fmt.Println("Get Post: This is post:", post)
		if err != nil {
			log.Printf("Post: getPost: %v", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		comments, err := h.services.GetAllComments(id)
		if err != nil {
			log.Println("Get Post: GetAllComments : ", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		info := models.Info{
			User:     user,
			Post:     post,
			Comments: comments,
		}
		ts, err := template.ParseFiles("./ui/templates/post.html")
		if err != nil {
			log.Printf("Get Post: Execute:%v", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
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

func (h *Handler) postLike(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/post-like/" {
	// 	h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	// 	return
	// }
	if r.Method == "GET" {
		user, ok := r.Context().Value(key).(models.User)
		if !ok {
			h.errorHandler(w, http.StatusInternalServerError, "Unauthorized")
			return
		}
		id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post-like/"))
		newPostLike := models.PostLike{
			UserID: user.ID,
			PostID: id,
			Status: 1,
		}
		if err := h.services.CreateLikePost(newPostLike); err != nil {
			log.Printf("Post: CreateLikePost: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		Idpost := strconv.Itoa(id)
		http.Redirect(w, r, "/get-post/"+Idpost, 302)
	}
}
