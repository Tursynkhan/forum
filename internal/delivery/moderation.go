package delivery

import (
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) createCategory(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create/category/" {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if user.RoleID != 3 {
		h.errorHandler(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	if err := r.ParseForm(); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Printf("%+v\n", r.Form)
	nameofcategory, ok := r.Form["createcategory"]
	if !ok {
		h.errorHandler(w, http.StatusInternalServerError, "createcategory not found")
		return
	}
	if err := h.services.CreateCategory(nameofcategory[0]); err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, "/profile/"+user.Username, http.StatusSeeOther)
}

func (h *Handler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete/category/" {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if user.RoleID != 3 {
		h.errorHandler(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	if err := r.ParseForm(); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
	categoryId, ok := r.Form["categories"]
	if !ok {
		h.errorHandler(w, http.StatusInternalServerError, "categories not found")
		return
	}
	catId, err := strconv.Atoi(strings.ReplaceAll(categoryId[0], " ", ""))
	if err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	fmt.Printf("%+v\n", r.Form)

	if err := h.services.DeleteCategoryById(catId); err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, "/profile/"+user.Username, http.StatusSeeOther)
}

func (h *Handler) promoteUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/promote/user/" {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if user.RoleID != 3 {
		h.errorHandler(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	if err := r.ParseForm(); err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Printf("%+v\n", r.Form)
	username, ok := r.Form["username"]
	if !ok {
		h.errorHandler(w, http.StatusInternalServerError, "username not found")
		return
	}
	rolesId, ok := r.Form["role"]
	if !ok {
		h.errorHandler(w, http.StatusInternalServerError, "role not found")
		return
	}
	roleId, err := strconv.Atoi(strings.ReplaceAll(rolesId[0], " ", ""))
	if err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if err := h.services.UpdateUserRole(username[0], roleId); err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, "/profile/"+user.Username, http.StatusSeeOther)

}
func (h *Handler) approve(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/approved/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if user.RoleID != 3 {
		h.errorHandler(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	if err := h.services.ApprovedPost(postId); err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, "/profile/"+user.Username, http.StatusSeeOther)

}
func (h *Handler) decline(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	user := r.Context().Value(key).(models.User)
	if user == (models.User{}) {
		h.errorHandler(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/post/declined/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if user.RoleID != 3 {
		h.errorHandler(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	if err := h.services.DeclinePost(postId); err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, "/profile/"+user.Username, http.StatusSeeOther)
}

func (h *Handler) reportOfPost(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		h.errorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	postId, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/posts/report/"))
	if err != nil {
		h.errorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	err = r.ParseForm()
	if err != nil {
		h.errorHandler(w, http.StatusInternalServerError, err.Error())
		return
	}
	nameReport, ok := r.Form["role"]
	if !ok {
		h.errorHandler(w, http.StatusInternalServerError, "role field not found")
		return
	}
	fmt.Println(nameReport)
	if err := h.services.UpdateReportOfPost(postId, nameReport[0]); err != nil {
		log.Println(err)
		h.errorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, "/post/"+strconv.Itoa(postId), http.StatusSeeOther)
}
