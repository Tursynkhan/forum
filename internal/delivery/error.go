package delivery

import (
	"html/template"
	"net/http"
)

type errStatus struct {
	StatusCode   int
	StatusString string
}

func (h *Handler) ErrorHandler(w http.ResponseWriter, r *http.Request, status errStatus) {
	w.WriteHeader(status.StatusCode)
	file := "./ui/template/error.html"
	ts, err := template.ParseFiles(file)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, status)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
