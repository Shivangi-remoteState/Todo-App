package dbHelper

import (
	"errors"
	"fmt"
	"mytodoApp/database"
	"mytodoApp/models"
	"strings"
	"time"
)

func CreateTodo(userID, name, description string, expiry_at time.Time) (models.Todos, error) {
	query := `
	INSERT INTO todos (user_id, name, description, expiry_At)
	VALUES ($1, $2, $3, $4)
	RETURNING id,
	          user_id,
	          name,
	          description,
	          completed_at,
	          expiry_at,
	          created_at;`
	var todo models.Todos
	err := database.DB.Get(&todo, query, userID, name, description, expiry_at)
	if err != nil {

		return models.Todos{}, err
	}

	return todo, nil
}

func GetTodos(userID string, search string, status string, limit, offset int) ([]models.Todos, error) {
	query := `SELECT id,
                    user_id,
                    name,
                    description,
                    completed_at,
                    expiry_at,
                    created_at
            FROM todos
            WHERE user_id = $1 
            AND archived_at is null`

	args := []interface{}{userID}
	i := 2

	//complete / incomplete / pending
	switch status {
	case "completed":
		query += " AND complete_at = true"

	case "pending":
		query += " AND complete_at = false AND expiry_at > NOW()"

	case "incomplete":
		query += " AND complete_at = false AND expiry_at <= NOW()"
	}
	//AND (search = "" or name ILIKE $2 )
	//search
	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d OR description ILIKE $%d", i, i)
		args = append(args, "%"+search+"%")
		i++
	}
	query += fmt.Sprintf(" ORDER BY expiry_at LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	todos := make([]models.Todos, 0)
	err := database.DB.Select(&todos, query, args...)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func GetTodoByID(todoID, userID string) (*models.Todos, error) {
	query := `SELECT id,
                     user_id,
                     name,
  				     description,
					 complete_at,
					 expiry_at,
 					 created_at
             FROM todos
             WHERE id = $1 AND user_id = $2
             AND archived_at is null`

	var todo models.Todos

	err := database.DB.Get(&todo, query, todoID, userID)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func UpdateTodo(req models.UpdateTodo, todoID, userID string) error {
	//query := `UPDATE todos
	//          SET name = $1,
	//              description =$2,
	//              complete = $3,
	//              expiry_at = $4
	//          WHERE id = $5
	//          AND  user_id = $6
	//          AND archived_at is null`

	query := "UPDATE todos SET "
	args := []interface{}{}
	i := 1

	if req.Name != nil {
		query += fmt.Sprintf("name = $%d, ", i)
		args = append(args, *req.Name)
		i++
	}
	if req.Description != nil {
		query += fmt.Sprintf("description = $%d, ", i)
		args = append(args, *req.Description)
		i++
	}
	if req.Complete != nil {
		query += fmt.Sprintf("complete = $%d, ", i)
		args = append(args, *req.Complete)
		i++
	}
	if req.ExpiryAt != nil {
		query += fmt.Sprintf("expiry_at = $%d, ", i)
		args = append(args, *req.ExpiryAt)
		i++
	}

	if len(args) == 0 {
		return errors.New("no fields to update")
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d AND archived_at is null", i, i+1)
	args = append(args, todoID, userID)

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("todo not found")
	}
	return nil
}

func DeleteTodoByID(todoID, userID string) error {
	query := `UPDATE todos 
              SET archived_at = NOW()
              WHERE id = $1 
              AND user_id = $2
              AND archived_at is null`

	result, err := database.DB.Exec(query, todoID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("todo not found")
	}
	return nil
}
