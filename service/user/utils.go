package user

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Mattcazz/Peer-Presure.git/service/auth"
	"github.com/Mattcazz/Peer-Presure.git/types"
	"github.com/Mattcazz/Peer-Presure.git/utils"
)

func loginUser(w http.ResponseWriter, userID int, username string) error {
	secret := os.Getenv("JWT_SECRET")

	token, err := auth.CreateJWT([]byte(secret), userID, username)

	if err != nil {
		return fmt.Errorf("Unauthorized")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	})

	return nil
}

func paginateFeed(h *Handler, userId, pageNumber int, username string) ([]*types.Post, *types.PaginationData, error) {

	totalPostsCount, err := h.postStore.GetPostsFromFriendsCount(userId)

	if err != nil {
		return nil, nil, err
	}

	posts, err := h.postStore.GetPostsFromFriends(userId, pageNumber, types.MaxPerPage)

	if err != nil {
		return nil, nil, err
	}

	pagination := utils.PreparePagination(pageNumber, totalPostsCount, fmt.Sprintf("/home/%s", username))

	return posts, pagination, nil
}
