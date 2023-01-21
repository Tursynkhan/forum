package delivery

import (
	"context"
	"errors"
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"time"
)

type contextKey string

const (
	key contextKey = "user"
)

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		token, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println(err)
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)
				return
			}
			h.errorHandler(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err = h.services.Autorization.ParseToken(token.Value)
		if err != nil {
			log.Printf("userIdentity : parse token %v", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		if user.Expiretime.Before(time.Now()) {
			if err := h.services.DeleteToken(token.Value); err != nil {
				h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), key, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
