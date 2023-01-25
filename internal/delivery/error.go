package delivery

import (
	"fmt"
	"net/http"
	"strings"
)

type data struct {
	Status  int
	Message string
}

func (h *Handler) errorHandler(w http.ResponseWriter, code int, errorText string) {
	w.WriteHeader(code)
	d := data{
		Status:  code,
		Message: http.StatusText(code),
	}
	if d.Status != 500 {
		temp := strings.Split(errorText, ":")
		d.Message = temp[len(temp)-1]
	}
	if err := h.tmpl.ExecuteTemplate(w, "error.html", d); err != nil {
		fmt.Fprintf(w, "%d - %s\n", d.Status, d.Message)
	}
}
