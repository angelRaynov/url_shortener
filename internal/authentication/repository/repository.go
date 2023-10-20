package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"strings"
	"url_shortener/internal/model"
)

const (
	QueryUserExist  = `SELECT count(*) FROM user WHERE email = ?`
	QueryInsertUser = `INSERT INTO user (uid, username, salt, password, email) VALUES (?, ?, ?, ?, ?)`
	QueryPatchUser  = `UPDATE user SET username = ?, salt = ?, password = ?, email = ? WHERE uid = ?;`
	QueryGetUser    = `SELECT uid, username, password, salt, email FROM user WHERE username = ?`
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (ar *authRepository) UserExist(email string) bool {

	row := ar.db.QueryRow(QueryUserExist, email)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false
	}

	return count > 0

}

func (ar *authRepository) GetUserByUsername(username string) (*model.User, error) {

	row := ar.db.QueryRow(QueryGetUser, username)

	var u model.User
	err := row.Scan(&u.UID, &u.Username, &u.Password, &u.Salt, &u.Email)
	if err != nil {
		return &model.User{}, err
	}

	return &u, nil

}

func (ar *authRepository) StoreUser(registerRequest *model.User, salt, hashedPass string) error {
	uid := uuid.New()

	// Store the salt and hashed password in the database
	_, err := ar.db.Exec(QueryInsertUser, uid.String(), registerRequest.Username, salt, hashedPass, registerRequest.Email)
	if err != nil {
		return err
	}

	return nil

}

func (ar *authRepository) EditUser(u *model.User, salt, uid string) error {
	q, args := buildPatchQuery(u, salt, uid)
	_, err := ar.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil

}

func buildPatchQuery(u *model.User, salt, uid string) (string, []interface{}) {
	q := "UPDATE user SET "
	var params []string
	var args []interface{}
	if u.Username != "" {
		args = append(args, u.Username)
		params = append(params, "username = ? ")
	}
	if u.Email != "" {
		args = append(args, u.Email)
		params = append(params, "email = ? ")
	}
	if u.Password != "" {
		args = append(args, u.Password, salt)
		params = append(params, "password = ?, salt = ? ")
	}

	args = append(args, uid)

	q += strings.Join(params, ", ")
	q += "WHERE uid = ?;"

	return q, args

}
