package users

import (
	"github.com/AJackTi/bookstore_users-api/datasource/mysql/users_db"
	"github.com/AJackTi/bookstore_users-api/utils/date_utils"
	"github.com/AJackTi/bookstore_users-api/utils/errors"
	"github.com/AJackTi/bookstore_users-api/utils/mysql_utils"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created FROM users where id = ?"
	queryUpdateUser = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?"
	queryDeleteUser = "DELETE FROM users WHERE id = ?"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.ID)
	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		return mysql_utils.ParseError(err)
	}

	userID, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
	}

	user.ID = userID
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)	
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.ID); err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}