package delivery

import (
	"net/http"

	"forum/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.userIdentity(h.home))
	mux.HandleFunc("/auth/signup", h.signUp)
	mux.HandleFunc("/auth/signin", h.signIn)
	mux.HandleFunc("/auth/logout", h.logout)

	mux.HandleFunc("/create-post", h.userIdentity(h.createPost))
	// mux.HandleFunc("/delete-post", h.userIdentity(h.deletePost))
	// mux.HandleFunc("/update-post", h.userIdentity(h.updatePost))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
