package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/forms"
	"forum/internal/models"
	"forum/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NewForms struct {
	Form     forms.Form
	Category []models.Category
}

const MAX_UPLOAD_SIZE = 20 * 1024 * 1024

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)

	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if r.URL.Path != "/create-post" {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	switch r.Method {

	case http.MethodGet:
		categories, err := h.services.GetAllCategories()
		if err != nil {
			log.Println("home page : get all categories", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		Form := NewForms{
			Form:     *forms.New(nil),
			Category: categories,
		}
		if err := h.tmpl.ExecuteTemplate(w, "createPost.html", Form); err != nil {
			log.Println(err.Error())
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
		if err != nil {
			log.Println("error parse form :", err)
			h.errorHandler(w, http.StatusBadRequest, "The uploaded file is too big.")
			return
		}
		title, ok := r.Form["title"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "title field not found")
			return
		}
		content, ok := r.Form["content"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "content field not found")
			return
		}
		categories, ok := r.Form["categories"]
		if !ok {
			h.errorHandler(w, http.StatusBadRequest, "categories field not found")
			return
		}

		files := r.MultipartForm.File["image"]
		if err != nil {
			fmt.Println("createPost: file-header: ", err)
			h.errorHandler(w, http.StatusBadRequest, err.Error())
			return
		}

		newPost := models.Post{
			UserID:  user.ID,
			Title:   title[0],
			Content: content[0],
			Files:   files,
			Created: time.Now().Format("2006-01-02 15:04:05"),
		}
		postId, err := h.services.Post.CreatePost(newPost)
		if err != nil {
			form := forms.New(r.PostForm)
			form.Required("title", "content", "categories")
			log.Printf("Post: Create Post: %v\n", err)
			if errors.Is(err, service.ErrPostTitleLen) {
				form.MaxLength("title", 100)
				w.WriteHeader(http.StatusBadRequest)
			} else if errors.Is(err, service.ErrPostContentLen) {
				form.MaxLength("content", 1500)
				w.WriteHeader(http.StatusBadRequest)
			} else if errors.Is(err, service.ErrInvalidPost) {
				form.Required("title", "content")
				w.WriteHeader(http.StatusBadRequest)
			} else if errors.Is(err, service.ErrInvalidType) {
				log.Println(err)
				h.errorHandler(w, http.StatusBadRequest, err.Error())
				return
			} else {
				log.Println(err)
				h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			if !form.Valid() {
				categories, err := h.services.GetAllCategories()
				if err != nil {
					log.Println("home page : get all categories", err)
					h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
					return
				}
				Form := NewForms{
					Form:     *form,
					Category: categories,
				}
				if err := h.tmpl.ExecuteTemplate(w, "createPost.html", Form); err != nil {
					log.Println(err.Error())
					h.errorHandler(w, http.StatusInternalServerError, err.Error())
					return
				}
				return
			}
		}
		if err = h.services.Post.CreatePostCategory(postId, categories); err != nil {
			log.Printf("Post: Create PostCategory : %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		log.Println("Create Post: Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return

	}
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		user := r.Context().Value(key).(models.User)
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/"))
		if err != nil {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		post, err := h.services.Post.GetPost(id)
		if err != nil {
			log.Printf("Post: getPost: %v", err)
			if errors.Is(err, sql.ErrNoRows) {
				h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		comments, err := h.services.GetAllComments(id)
		if err != nil {
			log.Println("Get Post: GetAllComments : ", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		for i := 0; i < len(comments); i++ {
			comments[i].Likes, err = h.services.GetCommentLikesByCommentID(comments[i].ID)
			if err != nil {
				log.Println("Get Post: GetCommentLikesByCommentID : ", err)
				h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			comments[i].Dislikes, err = h.services.GetCommentDislikesByCommentID(comments[i].ID)
			if err != nil {
				log.Println("Get Post: GetCommentLikesByCommentID : ", err)
				h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		likesPost, err := h.services.GetAllLikesByPostId(id)
		if err != nil {
			log.Println("Get Post: GetAllLikesByPostId : ", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		dislikesPost, err := h.services.GetAllDislikesByPostId(id)
		if err != nil {
			log.Println("Get Post: GetAllDislikesByPostId : ", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		newPostLike := models.PostLike{
			Likes:    likesPost,
			Dislikes: dislikesPost,
		}
		notifications, err := h.services.GetAllNotification(user)
		if err != nil {
			log.Println(err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		info := models.Info{
			User:          user,
			Post:          post,
			Comments:      comments,
			PostLike:      newPostLike,
			Notifications: notifications,
		}
		if err := h.tmpl.ExecuteTemplate(w, "post.html", info); err != nil {
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	default:
		log.Println("Get Post: Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))

	}
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if r.Method != http.MethodGet {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	postId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/delete/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, err.Error())
		return
	}
	post, err := h.services.Post.GetPost(postId)
	if err != nil {
		log.Printf("Post: getPost: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if err := h.services.DeletePost(post, user); err != nil {
		if errors.Is(err, service.ErrInvalidUser) {
			h.errorHandler(w, http.StatusBadRequest, err.Error())
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) editPost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	postId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/edit/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, err.Error())
		return
	}
	oldPost, err := h.services.GetPost(postId)
	if err != nil {
		log.Printf("Post: getPost: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
	title, ok := r.Form["update__title"]
	if !ok {
		h.errorHandler(w, http.StatusBadRequest, "title field not found")
		return
	}
	content, ok := r.Form["update__content"]
	if !ok {
		h.errorHandler(w, http.StatusBadRequest, "content field not found")
		return
	}
	categories, ok := r.Form["update__categories"]
	if !ok {
		h.errorHandler(w, http.StatusBadRequest, "categories field not found")
		return
	}
	newPost := models.Post{
		Title:   title[0],
		Content: content[0],
	}

	err := h.services.EditPost(oldPost, newPost, user)
	if err != nil {
	}
}

func (h *Handler) postLike(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post-like/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, err.Error())
		return
	}
	newPostLike := models.PostLike{
		UserID: user.ID,
		PostID: id,
		Status: 1,
	}
	post, err := h.services.Post.GetPost(id)
	if err != nil {
		log.Printf("Post: getPost: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if err := h.services.CreateLikePost(newPostLike); err != nil {
		log.Printf("Post: CreateLikePost: %v\n", err)
		if errors.Is(err, service.ErrPostNotexist) {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if post.Author != user.Username {
		newNotification := models.Notification{
			From:      user.Username,
			To:        post.Author,
			Content:   "liked your post",
			PostId:    post.ID,
			TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
			IsRead:    0,
		}
		if err := h.services.AddNewNotification(newNotification); err != nil {
			log.Println(err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	Idpost := strconv.Itoa(id)
	http.Redirect(w, r, "/post/"+Idpost, http.StatusSeeOther)
}

func (h *Handler) postDislike(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	user, ok := r.Context().Value(key).(models.User)
	if !ok {
		h.errorHandler(w, http.StatusInternalServerError, "Unauthorized")
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post-dislike/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, err.Error())
		return
	}
	newPostLike := models.PostLike{
		UserID: user.ID,
		PostID: id,
		Status: -1,
	}
	if err := h.services.CreateDisLikePost(newPostLike); err != nil {
		log.Printf("Post: CreateLikePost: %v\n", err)
		if errors.Is(err, service.ErrPostNotexist) {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	post, err := h.services.Post.GetPost(id)
	if err != nil {
		log.Printf("Post: getPost: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if post.Author != user.Username {
		newNotification := models.Notification{
			From:      user.Username,
			To:        post.Author,
			Content:   "disliked your post",
			PostId:    post.ID,
			TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
			IsRead:    0,
		}
		if err := h.services.AddNewNotification(newNotification); err != nil {
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	Idpost := strconv.Itoa(id)
	http.Redirect(w, r, "/post/"+Idpost, http.StatusSeeOther)
}
