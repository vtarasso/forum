package internal

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) SignUp(f *userSignupForm) error {
	user, err := m.UserByEmail(f.Email)
	if err != nil && !errors.Is(err, ErrNoRecord) {
		return err
	}
	if user.Email == f.Email {
		return ErrDuplicateEmail
	}
	user, err = m.UserByName(f.Name)
	if err != nil && !errors.Is(err, ErrNoRecord) {
		return err
	}
	if user.Name == f.Name {
		return ErrDuplicateName
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(f.Password), 12)
	if err != nil {
		return err
	}
	f.Password = string(hashedPassword)
	err = m.InsertUser(f)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) InsertUser(f *userSignupForm) error {
	query := `INSERT INTO users (name, email, hashed_password, created)
VALUES(?, ?, ?, datetime('now', '+6 hours'));`
	if _, err := m.DB.Exec(query, f.Name, f.Email, f.Password); err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UserByEmail(email string) (User, error) {
	query := `SELECT * FROM users
WHERE ? = email`
	var user User
	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Token, &user.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrNoRecord
	}
	return user, err
}

func (m *UserModel) UserByName(email string) (User, error) {
	query := `SELECT * FROM users
WHERE ? = name`
	var user User
	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Token, &user.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrNoRecord
	}
	return user, err
}

func (m *UserModel) GetUser(token string) (*User, error) {
	stmt := `SELECT id, name, email, hashed_password, created FROM users WHERE token = ?`
	row := m.DB.QueryRow(stmt, token)
	u := &User{}
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Если все хорошо, возвращается объект Snippet.
	return u, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) SetToken(id int, token string) error {
	query := `UPDATE users
SET token = ?, created = DATETIME('now', '+8 hours')
WHERE ? = id`
	_, err := m.DB.Exec(query, token, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) RemoveToken(token string) error {
	query := `UPDATE users
SET token = NULL
WHERE token = ?`
	_, err := m.DB.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) UserbyToken(sessionID string) (*User, error) {
	query := `SELECT * FROM users WHERE token = ?`
	user := &User{}
	err := m.DB.QueryRow(query, sessionID).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Token, &user.Created)
	if err != nil {
		return nil, err
	}
	return user, nil
}
