package post

import (
	"fmt"
	"math"

	"github.com/Mattcazz/Peer-Presure.git/types"
	"github.com/Mattcazz/Peer-Presure.git/utils"
)

func paginatePosts(h *Handler, userId, pageNumber int, username string) ([]*types.Post, *types.PaginationData, error) {

	totalPostsCount, err := h.postStore.GetPostsFromUserCount(userId)

	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(totalPostsCount) / float64(types.MaxPerPage)))

	if totalPages == 0 {
		totalPages = 1
	}

	posts, err := h.postStore.GetPostsFromUser(pageNumber, types.MaxPerPage, username)

	if err != nil {
		return nil, nil, err
	}

	var pagination *types.PaginationData

	if totalPostsCount > types.MaxPerPage {
		baseURL := fmt.Sprintf("/%s/posts", username)
		pagination = utils.PreparePagination(pageNumber, totalPages, baseURL)
	}

	return posts, pagination, nil
}
