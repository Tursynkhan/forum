package delivery

import (
	"forum/internal/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.userIdentity(h.home))
	mux.HandleFunc("/auth/signup", h.signUp)
	mux.HandleFunc("/auth/signin", h.signIn)
	mux.HandleFunc("/auth/logout", h.logout)

	mux.HandleFunc("/create-post", h.userIdentity(h.createPost))
	mux.HandleFunc("/get-post/", h.userIdentity(h.getPost))
	mux.HandleFunc("/post-like/", h.userIdentity(h.postLike))
	mux.HandleFunc("/post-dislike/", h.userIdentity(h.postDislike))

	mux.HandleFunc("/create-comment", h.userIdentity(h.createComment))
	mux.HandleFunc("/comment-like/", h.userIdentity(h.commentLike))
	mux.HandleFunc("/comment-dislike/", h.userIdentity(h.commentDislike))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return h.recoverPanic(h.logRequest(h.secureHeaders(mux)))
}
