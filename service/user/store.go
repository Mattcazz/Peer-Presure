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

func (s *Store) CreateUser(user types.User) error {

	query := `INSERT INTO users 
			(user_name, email, password, created_at)
			VALUES ($1, $2, $3, $4)`

	_, err := s.db.Query(query,
		user.UserName,
		user.Email,
		user.Password,
		user.CreatedAt)

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

// Private function that returns a user given a row to scan
func scanUserRow(row *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.CreatedAt)

	return user, err
}
