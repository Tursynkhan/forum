package delivery

import (
	"forum/internal/service"
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl     *template.Template
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		tmpl:     template.Must(template.ParseGlob("./ui/templates/*.html")),
		services: services,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	m := new(middleware)
	m.addMidlleware(h.userIdentity)
	m.addMidlleware(h.logRequest)
	m.addMidlleware(h.secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc("/", m.chain(h.home))
	mux.HandleFunc("/auth/signup", m.chain(h.signUp))
	mux.HandleFunc("/auth/signin", m.chain(h.signIn))
	mux.HandleFunc("/auth/logout", m.chain(h.logout))
	mux.HandleFunc("/profile/", m.chain(h.profilePage))

	mux.HandleFunc("/create-post", m.chain(h.createPost))
	mux.HandleFunc("/post/", m.chain(h.getPost))
	mux.HandleFunc("/post/delete/", m.chain(h.deletePost))
	mux.HandleFunc("/post/edit/", m.chain(h.editPost))
	mux.HandleFunc("/post-like/", m.chain(h.postLike))
	mux.HandleFunc("/post-dislike/", m.chain(h.postDislike))

	mux.HandleFunc("/create-comment", m.chain(h.createComment))
	mux.HandleFunc("/comment/delete/", m.chain(h.deleteComment))
	mux.HandleFunc("/comment/edit/", m.chain(h.editComment))
	mux.HandleFunc("/comment-like/", m.chain(h.commentLike))
	mux.HandleFunc("/comment-dislike/", m.chain(h.commentDislike))

	mux.HandleFunc("/create/category/", m.chain(h.createCategory))
	mux.HandleFunc("/delete/category/", m.chain(h.deleteCategory))
	mux.HandleFunc("/promote/user/", m.chain(h.promoteUser))
	mux.HandleFunc("/post/approved/", m.chain(h.approve))
	mux.HandleFunc("/post/declined/", m.chain(h.decline))
	mux.HandleFunc("/posts/report/",m.chain(h.reportOfPost))
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
