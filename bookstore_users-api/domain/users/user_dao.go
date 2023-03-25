package users

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/AJackTi/bookstore_users-api/datasource/mysql/users_db"
	"github.com/AJackTi/bookstore_utils-go/logger"
	"github.com/AJackTi/bookstore_utils-go/rest_errors"
)

const (
	indexUniqueEmail            = "email_UNIQUE"
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status FROM users where id = ?"
	queryUpdateUser             = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?"
	queryDeleteUser             = "DELETE FROM users WHERE id = ?"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email = ? AND password = ? AND status = ?"
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.ID)
	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when trying to get user by id", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}

	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}
	defer stmt.Close()

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if err != nil {
		logger.Error("error when trying to save user", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}

	userID, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}

	user.ID = userID
	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}

	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.ID); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}

	return nil
}

func (user *User) Search(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when scan user row into user struct", err)
			return nil, rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, rest_errors.New(http.StatusNotFound, errors.New(fmt.Sprintf("no users matching status %s", status)))
	}

	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		if err == sql.ErrNoRows {
			return rest_errors.New(http.StatusNotFound, errors.New("invalid user credentials"))
		}
		logger.Error("error when trying to get user by email and password", err)
		return rest_errors.New(http.StatusInternalServerError, errors.New("database error"))
	}

	return nil
}
