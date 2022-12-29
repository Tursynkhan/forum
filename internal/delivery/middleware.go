package delivery

import (
	"context"
	"fmt"
	"net/http"
)

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("session_token")
		if token.Value == "" {
			next.ServeHTTP(w, r)
		}
		fmt.Println(token)
		user, err := h.services.Autorization.ParseToken(token.Value)
		if err != nil {
			h.errorHandler(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println(user)
		ctx := context.WithValue(r.Context(), "myvalues", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
