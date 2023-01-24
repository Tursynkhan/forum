package delivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
)

type contextKey string

const key contextKey = "user"

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			if errors.Is(err, sql.ErrNoRows) {
				next.ServeHTTP(w, r)
				return
			}
			log.Printf("userIdentity : parse token %v", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := h.services.Session.CompareExpirationTime(); err != nil {
			log.Println("middleware : userIdentity: ", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), key, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XXS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
