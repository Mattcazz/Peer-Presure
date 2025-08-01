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
	query := `INSERT INTO posts (user_id, username,title, text, img_url, public,created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING *`

	row, err := s.db.Query(query,
		p.UserId,
		p.Username,
		p.Title,
		p.Text,
		p.ImgURL,
		p.Public,
		p.CreatedAt)

	if err != nil {
		return nil, err
	}

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
func (s *Store) GetPostsFromUser(page, maxPerPage int, username string) ([]*types.Post, error) {

	offset := (page - 1) * maxPerPage

	query := `SELECT * FROM posts WHERE username = $1 ORDER BY created_at DESC OFFSET $2 LIMIT $3`

	var posts []*types.Post
	var post *types.Post

	rows, err := s.db.Query(query, username, offset, maxPerPage)

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

func (s *Store) GetLastTenPosts() ([]*types.Post, error) {
	query := ` SELECT * FROM posts 
			WHERE public = true
			ORDER BY created_at DESC
			LIMIT 10 `

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}
	var posts []*types.Post
	var post *types.Post

	for rows.Next() {
		post, err = scanPostRow(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Store) EditPost(post *types.Post) error {
	query := ` UPDATE posts
			   SET title = $1, text = $2, img_url = $3, public = $4
			   WHERE id = $5`
	_, err := s.db.Query(query, post.Title, post.Text, post.ImgURL, post.Public, post.ID)

	return err
}

func (s *Store) GetPostsFromFriends(userId, page, maxPerPage int) ([]*types.Post, error) {

	offset := (page - 1) * maxPerPage

	query := `SELECT p.* FROM posts p 
			  JOIN users u ON p.user_id = u.id
			  JOIN friends f   ON (
      			(f.user_id1 = $1 AND p.user_id = f.user_id2) OR (f.user_id2 = $1 AND p.user_id = f.user_id1)
			  )
			  WHERE p.user_id != $1 AND f.status = 'confirmed'
			  ORDER BY p.created_at DESC
			  OFFSET $2
			  LIMIT $3`

	rows, err := s.db.Query(query, userId, offset, maxPerPage)

	if err != nil {
		return nil, err
	}

	var posts []*types.Post
	var post *types.Post

	for rows.Next() {
		post, err = scanPostRow(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil

}

func (s *Store) GetPostsFromFriendsCount(userId int) (int, error) {
	var totalCount int

	query := `SELECT COUNT(*) FROM posts p 
			  JOIN users u ON p.user_id = u.id
			  JOIN friends f   ON (
      			(f.user_id1 = $1 AND p.user_id = f.user_id2) OR (f.user_id2 = $1 AND p.user_id = f.user_id1)
			  )
			  WHERE p.user_id != $1 AND f.status = 'confirmed'`

	err := s.db.QueryRow(query, userId).Scan(&totalCount)

	return totalCount, err
}

func (s *Store) GetPostsFromUserCount(userId int) (int, error) {
	var totalCount int

	query := `SELECT COUNT(*) FROM posts p 
			  WHERE user_id = $1`

	err := s.db.QueryRow(query, userId).Scan(&totalCount)

	return totalCount, err
}

func scanPostRow(r *sql.Rows) (*types.Post, error) {
	post := new(types.Post)

	err := r.Scan(
		&post.ID,
		&post.UserId,
		&post.Username,
		&post.Title,
		&post.Text,
		&post.ImgURL,
		&post.Public,
		&post.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return post, nil
}
