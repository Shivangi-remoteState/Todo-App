package dbHelper

import (
	"errors"
	"mytodoApp/database"
)

func IsUserExist(email string) (bool, error) {
	query := `SELECT count(*)>0
              FROM users 
              where email =TRIM(LOWER($1)) 
              and archived_At is null `

	var IsUserExist bool
	err := database.DB.Get(&IsUserExist, query, email)
	return IsUserExist, err
}

func CreateUser(name, email, password string) (string, error) {

	query := `INSERT INTO users (name, email, password)
			  VALUES ($1, TRIM(LOWER($2)), $3)
			  RETURNING id;`

	var userID string

	err := database.DB.Get(&userID, query, name, email, password)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func CreateUserSession(userID string) (string, error) {

	query := `INSERT INTO user_session (user_id)
			  VALUES ($1)
			  RETURNING id;`

	var sessionID string

	err := database.DB.Get(&sessionID, query, userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func GetUserBYEmail(email string) (string, string, string, error) {
	query := `SELECT id, password, role 
              FROM users where email = TRIM(LOWER($1))
              AND archived_At is null
              AND suspended_at IS NULL`

	var result struct {
		ID       string `db:"id"`
		Password string `db:"password"`
		Role     string  `db:"role"`
	}
	err := database.DB.Get(&result, query, email)
	if err != nil {
		return "", "", "", err
	}

	return result.ID, result.Password, result.Role,  nil
}

// logout
func ArchiveUserSession(token string) error {
	query := `UPDATE user_session
              set archived_at = NOW()
              where id = $1
              and archived_At is null;`

	result, err := database.DB.Exec(query, token)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("invalid session")
	}
	return nil
}

func GetUserIDBySessionID(token string) (string, error) {
	query := `SELECT user_id FROM user_session WHERE id = $1 and archived_AT is null`
	var userID string

	err := database.DB.Get(&userID, query, token)
	if err != nil {
		return "", errors.New("Invalid session")
	}
	return userID, nil
}
