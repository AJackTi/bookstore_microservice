package mysql_utils

import (
	"database/sql"

	"github.com/AJackTi/bookstore_users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if err == sql.ErrNoRows {
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerError("error parsing database response")
	}

	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError("invalid data")
	}
		
	return errors.NewInternalServerError("error processing request")
}