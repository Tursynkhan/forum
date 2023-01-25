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

type middleware struct {
	allMiddlewares []func(http.HandlerFunc) http.HandlerFunc
}

func (m *middleware) addMidlleware(middle func(http.HandlerFunc) http.HandlerFunc) {
	m.allMiddlewares = append(m.allMiddlewares, middle)
}

func (m *middleware) chain(router http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range m.allMiddlewares {
		router = middleware(router)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	}
}

type contextKey string

const key contextKey = "user"

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		token, err := r.Cookie("token")
		if err != nil {
			if err := h.services.DeleteTokenWhenExpireTime(); err != nil {
				h.errorHandler(w, http.StatusInternalServerError, err.Error())
				return
			}
			if errors.Is(err, http.ErrNoCookie) {
				fmt.Println("error", err)
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), key, models.User{})))
				return
			}
			h.errorHandler(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err = h.services.Autorization.ParseToken(token.Value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), key, models.User{})))
				return
			}
			log.Printf("userIdentity : parse token %v", err)
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), key, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) secureHeaders(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XXS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s\t[%s]\t%s%s", r.Proto, r.Method, r.Host, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
