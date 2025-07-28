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

func (s *Store) CreateComment(c *types.Comment) (*types.Comment, error) {
	query := `INSERT INTO comments (user_id, post_id, text, created_at)
	VALUES ($1, $2, $3, $4) RETURNING *`

	row, err := s.db.Query(query,
		c.UserID,
		c.PostID,
		c.Text,
		c.CreatedAt)

	if err != nil {
		return nil, err
	}

	if row.Next() {
		return scanCommentRow(row)
	}

	return nil, err
}

func (s *Store) DeleteComment(id int) error {
	query := `DELETE FROM comments WHERE ID = $1`

	_, err := s.db.Query(query, id)

	return err
}

func (s *Store) GetCommentsFromPost(post_id int) ([]*types.Comment, error) {
	rows, err := s.db.Query("SELECT * FROM comments WHERE post_id = $1", post_id)

	if err != nil {
		return nil, err
	}

	comments := []*types.Comment{}

	for rows.Next() {

		comment, err := scanCommentRow(rows)

		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil

}

func (s *Store) GetCommentsFromUser(user_id int) ([]*types.Comment, error) {
	rows, err := s.db.Query("SELECT * FROM comments WHERE user_id = $1", user_id)

	if err != nil {
		return nil, err
	}

	comments := []*types.Comment{}

	for rows.Next() {

		comment, err := scanCommentRow(rows)

		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func scanCommentRow(row *sql.Rows) (*types.Comment, error) {

	comment := new(types.Comment)

	err := row.Scan(
		&comment.ID,
		&comment.UserID,
		&comment.PostID,
		&comment.Text,
		&comment.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return comment, nil
}
