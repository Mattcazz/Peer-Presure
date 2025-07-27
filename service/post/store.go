package post

import (
	"database/sql"

	"github.com/Mattcazz/Peer-Presure.git/types"
)

type Store struct {
	db *sql.DB
}

// CreatePost implements types.PostStore.
func (s *Store) CreatePost(types.Post) error {
	panic("unimplemented")
}

// DeletePost implements types.PostStore.
func (s *Store) DeletePost(int) error {
	panic("unimplemented")
}

// GetPostById implements types.PostStore.
func (s *Store) GetPostById(int) (*types.Post, error) {
	panic("unimplemented")
}

// GetPostsFromUser implements types.PostStore.
func (s *Store) GetPostsFromUser(int) ([]*types.Post, error) {
	panic("unimplemented")
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func GetPostById(id int) (*types.Post, error) {
	return nil, nil
}

func GetPostsFromUser(user_id int) ([]*types.Post, error) {
	return nil, nil
}

func CreatePost(p types.Post) error {
	return nil
}

func DeletePost(post_id int) error {
	return nil
}
