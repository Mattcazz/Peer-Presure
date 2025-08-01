package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/Mattcazz/Peer-Presure.git/types"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {

	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func GetUserIdFromRequest(r *http.Request) (int, error) {

	userId, ok := r.Context().Value(types.CtxKeyUserID).(int)

	if !ok {
		return 0, fmt.Errorf("error: with jwt userID")
	}

	return userId, nil
}
func GetUsernameFromRequest(r *http.Request) (string, error) {

	username, ok := r.Context().Value(types.CtxKeyUsername).(string)

	if !ok {
		return "", fmt.Errorf("error: with jwt username")
	}

	return username, nil
}
func GetIdFromURL(format string, r *http.Request) (int, error) {

	vars := mux.Vars(r)
	id := vars[format]
	postID, err := strconv.Atoi(id)

	if err != nil {
		return 0, err
	}

	return postID, nil
}

func PreparePagination(currentPage, totalPostsCount int, baseUrl string) *types.PaginationData {

	totalPages := int(math.Ceil(float64(totalPostsCount) / float64(types.MaxPerPage)))

	if totalPages == 0 {
		totalPages = 1
	}

	var pagination *types.PaginationData

	if totalPostsCount > types.MaxPerPage {
		pagination = paginate(currentPage, totalPages, baseUrl)
	}

	return pagination
}

func paginate(currentPage, totalPages int, baseUrl string) *types.PaginationData {
	var pageNumbers []int
	for i := 1; i <= totalPages; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	return &types.PaginationData{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		BaseURL:     baseUrl,
		PageNumbers: pageNumbers,
		PrevPage:    currentPage - 1,
		NextPage:    currentPage + 1,
		HasPrev:     currentPage > 1,
		HasNext:     currentPage < totalPages,
	}
}
