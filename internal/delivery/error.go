package delivery

import (
	"net/http"
	"strings"
	"text/template"
)

type data struct {
	Status  int
	Message string
	ErrText string
}

func (h *Handler) errorHandler(w http.ResponseWriter, code int, errorText string) {
	w.WriteHeader(code)
	d := data{
		Status:  code,
		Message: http.StatusText(code),
		ErrText: errorText,
	}
	if d.Status != http.StatusInternalServerError {
		temp := strings.Split(errorText, ":")
		d.ErrText = temp[len(temp)-1]
	}
	ts, err := template.ParseFiles("./ui/templates/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, d)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
