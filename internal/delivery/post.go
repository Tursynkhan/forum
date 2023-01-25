package delivery

import (
	"database/sql"
	"errors"
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

	if r.Method == "GET" {

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
		categories := r.Form["categories"]
		newPost := models.Post{
			UserID:  user.ID,
			Title:   title,
			Content: content,
			Created: time.Now().Format("2006-01-02 15:04:05"),
		}
		postId, err := h.services.Post.CreatePost(newPost)
		if err != nil {
			form := forms.New(r.PostForm)
			log.Printf("Post: Create Post: %v\n", err)
			if errors.Is(err, service.ErrPostTitleLen) {
				form.MaxLength("title", 100)
			} else if errors.Is(err, service.ErrPostContentLen) {
				form.MaxLength("content", 1500)
			} else if errors.Is(err, service.ErrInvalidPost) {
				form.Required("title", "content")
			} else {
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
	} else {
		log.Println("Create Post: Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) getPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		user := r.Context().Value(key).(models.User)

		id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/get-post/"))

		post, err := h.services.Post.GetPost(id)
		if err != nil {
			log.Printf("Post: getPost: %v", err)
			if errors.Is(err, sql.ErrNoRows) {
				h.errorHandler(w, http.StatusNotFound, err.Error())
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

		info := models.Info{
			User:     user,
			Post:     post,
			Comments: comments,
			PostLike: newPostLike,
		}
		if err := h.tmpl.ExecuteTemplate(w, "post.html", info); err != nil {
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	} else {
		log.Println("Get Post: Method not allowed")
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func (h *Handler) postLike(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if r.Method == "POST" {
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
		if err := h.services.CreateLikePost(newPostLike); err != nil {
			log.Printf("Post: CreateLikePost: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		Idpost := strconv.Itoa(id)
		http.Redirect(w, r, "/get-post/"+Idpost, http.StatusSeeOther)
	} else {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *Handler) postDislike(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if r.Method == "POST" {
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
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		Idpost := strconv.Itoa(id)
		http.Redirect(w, r, "/get-post/"+Idpost, http.StatusSeeOther)
	} else {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
