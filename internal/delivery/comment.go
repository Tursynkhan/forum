package delivery

import (
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user, ok := r.Context().Value(key).(models.User)
		if !ok {
			h.errorHandler(w, http.StatusInternalServerError, "Unauthorized")
			return
		}
		postId, _ := strconv.Atoi(r.PostFormValue("postId"))
		comment := r.PostFormValue("comment")

		newComment := models.Comment{
			Content: comment,
			UserID:  user.ID,
			PostID:  postId,
		}
		if err := h.services.CreateComment(newComment); err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		Idpost := strconv.Itoa(postId)
		http.Redirect(w, r, "/get-post/"+Idpost, 302)
	}
}

func (h *Handler) commentLike(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		user, ok := r.Context().Value(key).(models.User)
		if !ok {
			h.errorHandler(w, http.StatusInternalServerError, "Unauthorized")
			return
		}
		id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment-like/"))
		newCommentLike := models.CommentLike{
			UserID:    user.ID,
			CommentID: id,
			Status:    1,
		}
		comment, err := h.services.GetCommentById(id)
		if err != nil {
			log.Printf("Comment: CommentLike: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if err := h.services.CreateLikeComment(newCommentLike); err != nil {
			log.Printf("Comment: CommentLike: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/get-post/%d", comment.PostID), 302)
	}
}

func (h *Handler) commentDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		user, ok := r.Context().Value(key).(models.User)
		if !ok {
			h.errorHandler(w, http.StatusInternalServerError, "Unauthorized")
			return
		}
		id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment-dislike/"))
		newCommentLike := models.CommentLike{
			UserID:    user.ID,
			CommentID: id,
			Status:    -1,
		}
		comment, err := h.services.GetCommentById(id)
		if err != nil {
			log.Printf("Comment: CommentLike: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if err := h.services.CreateDisLikeComment(newCommentLike); err != nil {
			log.Printf("Comment: CommentDislike: %v\n", err)
			h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/get-post/%d", comment.PostID), 302)
	}
}
