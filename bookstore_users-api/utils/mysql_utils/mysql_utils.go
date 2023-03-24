package mysql_utils

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/AJackTi/bookstore_utils-go/rest_errors"
	"github.com/go-sql-driver/mysql"
)

func ParseError(err error) *rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if err == sql.ErrNoRows {
			return rest_errors.New(http.StatusNotFound, errors.New("no record matching given id"))
		}
		return rest_errors.New(http.StatusInternalServerError, errors.New("error parsing database response"))
	}

	switch sqlErr.Number {
	case 1062:
		return rest_errors.New(http.StatusBadRequest, errors.New("invalid data"))
	}

	return rest_errors.New(http.StatusInternalServerError, errors.New("error processing request"))
}
