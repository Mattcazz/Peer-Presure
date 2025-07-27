package comment

import (
	"database/sql"

	"github.com/Mattcazz/Peer-Presure.git/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateComment(*types.Comment) error {
	panic("unimplemented")
}

func (s *Store) DeleteComment(int) error {
	panic("unimplemented")
}

func (s *Store) GetCommentsFromPost(int) ([]*types.Comment, error) {
	panic("unimplemented")
}

func (s *Store) GetCommentsFromUser(int) ([]*types.Comment, error) {
	panic("unimplemented")
}
