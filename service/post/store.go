package post

import (
	"database/sql"
	"fmt"

	"github.com/Mattcazz/Peer-Presure.git/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreatePost implements types.PostStore.
func (s *Store) CreatePost(p types.Post) (*types.Post, error) {
	query := `INSERT INTO posts (user_id, title, text, img_url, public,created_at)
	VALUES ($1,$2,$3,$4,$5,$6) RETURNING *`

	row, err := s.db.Query(query,
		p.UserId,
		p.Title,
		p.Text,
		p.ImgURL,
		p.Public,
		p.CreatedAt)

	if row.Next() {
		return scanPostRow(row)
	}

	return nil, err
}

// DeletePost implements types.PostStore.
func (s *Store) DeletePost(post_id int) error {
	query := `DELETE FROM posts WHERE ID = $1`

	_, err := s.db.Query(query, post_id)

	return err
}

// GetPostById implements types.PostStore.
func (s *Store) GetPostById(post_id int) (*types.Post, error) {
	query := `SELECT * FROM posts WHERE ID = $1`

	row, err := s.db.Query(query, post_id)

	if err != nil {
		return nil, err
	}

	if row.Next() {
		return scanPostRow(row)
	}

	return nil, fmt.Errorf("the search came up with no results")

}

// GetPostsFromUser implements types.PostStore.
func (s *Store) GetPostsFromUser(user_id int) ([]*types.Post, error) {
	query := `SELECT * FROM posts WHERE user_id = $1`

	var posts []*types.Post
	var post *types.Post

	rows, err := s.db.Query(query, user_id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		post, err = scanPostRow(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func scanPostRow(r *sql.Rows) (*types.Post, error) {
	post := new(types.Post)

	err := r.Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Text,
		&post.ImgURL,
		&post.Text,
		&post.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return post, nil
}
