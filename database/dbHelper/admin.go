package dbHelper

import (
	"errors"
	"fmt"
	"mytodoApp/database"
	"mytodoApp/models"
)

func GetAllTodos(search, status string) ([]models.Todos, error) {
	query := `SELECT id,
	                 user_id,
	                 name,
	                 description,
	                 completed_at,
	                 expiry_at,
	                 created_at
	          FROM todos
	          WHERE archived_at IS NULL`

	args := []interface{}{}
	i := 1

	// status filter
	switch status {
	case "completed":
		query += " AND completed_at = true"

	case "pending":
		query += " AND completed_at = false AND expiry_at > NOW()"

	case "incomplete":
		query += " AND completed_at = false AND expiry_at <= NOW()"
	}

	// search filter
	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", i, i)
		args = append(args, "%"+search+"%")
		i++
	}

	query += " ORDER BY expiry_at"

	todos := []models.Todos{}
	err := database.DB.Select(&todos, query, args...)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func GetAllUsers(search string) ([]models.User, error) {

	query := `SELECT id,
	                 name,
	                 email,
	                 role,
	                 created_at,
	                 suspended_at
	          FROM users
	          WHERE archived_at IS NULL`

	args := []interface{}{}
	i := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", i, i)
		args = append(args, "%"+search+"%")
		i++
	}

	query += " ORDER BY created_at DESC"

	users := []models.User{}
	err := database.DB.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// suspend user
func SuspendUser(userID string) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//suspend user
	_, err = tx.Exec(
		`UPDATE users SET suspended_at = NOW() where id = $1`, userID)
	if err != nil {
		return err
	}

	//	invalidate all session
	_, err = tx.Exec(`UPDATE  user_session SET archived_at = NOW() where user_id = $1`, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func UnsuspendUser(userID string) error {

	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`UPDATE users SET suspended_at =null WHERE id = $1`, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not suspended or not found")
	}
	return tx.Commit()
}
