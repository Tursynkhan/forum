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
		fmt.Println("middleware : ParseToken : user :")
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

func (h *Handler) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XXS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

// func (h *Handler) myMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !isAuthorized(r) {
// 			h.errorHandler(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }
func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// func (h *Handler) recoverPanic(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				w.Header().Set("Connection", "close")
// 				log.Printf("middleware : recoverPanic: ", err)
// 				h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

// 			}
// 		}()
// 		next.ServeHTTP(w, r)
// 	})
// }
