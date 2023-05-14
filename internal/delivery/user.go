package delivery

import (
	"errors"
	"forum/internal/models"
	"forum/internal/service"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) profilePage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(key).(models.User)
	if r.Method != http.MethodGet {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	userPage, err := h.services.GetProfileByUsername(username)
	if err != nil {
		log.Println("profile : GetProfileByUsername", err)
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}
	posts, err := h.services.GetPostByUsername(userPage.Username, r.URL.Query())
	if err != nil {
		log.Println("user : GetPostByUsername :", err)
		if errors.Is(err, service.ErrInvalidQuery) {
			h.errorHandler(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	notifications, err := h.services.GetAllNotification(user)
	if err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	info := models.Info{
		User:          user,
		ProfileUser:   userPage,
		Posts:         posts,
		Notifications: notifications,
	}
	if err := h.tmpl.ExecuteTemplate(w, "profile.html", info); err != nil {
		log.Println("user : Executetemplate : ", err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}
