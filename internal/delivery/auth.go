package delivery

import (
	"net/http"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	// var input models.User
	// id,err:=h.services.Autorization.CreateUser(input)
	// if err!=nil{
	// 	http.Error(w,"Internal server error",500)
	// 	return
	// }
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
}
