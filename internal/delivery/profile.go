package delivery

import (
	"forum/internal/models"
	"net/http"
)

func (h *Handler) profilePage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if r.Method != http.MethodGet {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
	// username := strings.TrimPrefix(r.URL.Path, "/profile")
	// userPage, err := h.services.GetUserByUsername(username)
	// if err != nil {
	// 	log.Println("profile : GetUserByUsername", err)
	// 	h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	// }
	// posts, err := h.services.GetPostByUsername(username, r.URL.Query())
	// if err != nil {
	// 	h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	return
	// }
	info := models.Info{
		User: user,
	}
	if err := h.tmpl.ExecuteTemplate(w, "profile.html", info); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}
