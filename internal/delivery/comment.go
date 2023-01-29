package delivery

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postId, err := strconv.Atoi(r.PostFormValue("postId"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	comment, ok := r.Form["comment"]
	if !ok {
		h.errorHandler(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	if len(comment[0]) == 0 {
		h.errorHandler(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	newComment := models.Comment{
		Content: comment[0],
		UserID:  user.ID,
		PostID:  postId,
	}
	if err := h.services.CreateComment(newComment); err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrPostNotexist) {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
	Idpost := strconv.Itoa(postId)
	http.Redirect(w, r, "/post/"+Idpost, http.StatusSeeOther)
}

func (h *Handler) commentLike(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, "can not like comment")
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment-like/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
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
		if errors.Is(err, service.ErrCommentNotExist) {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostID), http.StatusSeeOther)
}

func (h *Handler) commentDislike(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, "can not like comment")
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/comment-dislike/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
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
		if errors.Is(err, service.ErrCommentNotExist) {
			h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostID), http.StatusSeeOther)
}
