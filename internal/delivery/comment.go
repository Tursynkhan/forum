package delivery

import (
	"net/http"
	"strconv"

	"forum/internal/models"
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
		http.Redirect(w, r, "/get-post?id="+Idpost, 302)
	}
}
