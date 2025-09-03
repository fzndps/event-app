package database

import (
	"context"
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

type UserRepository interface {
	InsertUser(ctx context.Context, user User) User
}

func (m *UserModel) InsertUser(ctx context.Context, user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := "INSERT INTO users (email, name, password) VALUES (?, ?, ?)"

	result, err := m.DB.ExecContext(ctx, query, user.Email, user.Name, user.Password)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = int(id)

	return user, nil
}

func (m *UserModel) getUser(ctx context.Context, query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetUserById(id int) (*User, error) {
	query := "SELECT id, email, name, password FROM users WHERE id = ?"
	return m.getUser(context.Background(), query, id)
}

func (m *UserModel) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, email, name, password FROM users WHERE email = ?"
	return m.getUser(context.Background(), query, email)
}
