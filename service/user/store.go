package user

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

func (s *Store) GetUserByEmail(email string) (*types.User, error) {

	query := `SELECT * FROM users WHERE "email" = $1`

	row, err := s.db.Query(query, email)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		return scanUserRow(row)
	}

	return nil, fmt.Errorf("the search came up with no results")
}

func (s *Store) CreateUser(user *types.User) error {

	query := `INSERT INTO users 
			(email, username, password, created_at)			
			VALUES ($1, $2, $3, $4) RETURNING *`

	err := s.db.QueryRow(query,
		user.Email,
		user.UserName,
		user.Password,
		user.CreatedAt).Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.Password,
		&user.CreatedAt)

	return err
}

func (s *Store) GetUserById(id int) (*types.User, error) {

	query := `SELECT * FROM users WHERE id = $1`

	row, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		return scanUserRow(row)
	}

	return nil, fmt.Errorf("the search came up with no results")

}

func (s *Store) GetUserByUsername(username string) (*types.User, error) {
	query := `SELECT * FROM users WHERE username = $1`

	row, err := s.db.Query(query, username)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		return scanUserRow(row)
	}

	return nil, fmt.Errorf("the search came up with no results")
}

func (s *Store) GetUserFriends(userId int) ([]*types.User, error) {
	query := `SELECT u.* FROM users u
			  JOIN friends f ON 
  				(f.user_id1 = $1 AND u.id = f.user_id2) OR
  				(f.user_id2 = $1 AND u.id = f.user_id1)
			  ORDER BY u.username ASC`

	var users []*types.User
	var user *types.User

	rows, err := s.db.Query(query, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user, err = scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil

}

func (s *Store) CreateFriendRequest(id1, id2 int) error {
	_, err := s.db.Query("INSERT INTO friends (user_id1, user_id2) VALUES ($1, $2)", id1, id2)

	return err
}

func (s *Store) DeleteFriend(id1, id2 int) error {

	query := `DELETE FROM friends 
			  WHERE (user_id1 = $1 AND user_id2 = $2) OR (user_id1 = $2 AND user_id2 = $1)`
	_, err := s.db.Query(query, id1, id2)

	return err
}

func (s *Store) RespondFriendRequest(id1, id2 int, r string) error {
	query := `UPDATE friends SET status = $3
			  WHERE (user_id1 = $1 AND user_id2 = $2) OR (user_id1 = $2 AND user_id2 = $1) `

	_, err := s.db.Query(query, r, id1, id2)

	return err
}

// Private function that returns a user given a row to scan
func scanUserRow(row *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.Password,
		&user.CreatedAt)

	return user, err
}
