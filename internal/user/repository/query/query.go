package query

import (
	"fmt"
	"strings"

	s "github.com/core-go/sql"

	"go-service/internal/user/model"
)

func BuildQuery(filter *model.UserFilter) (string, []interface{}) {
	query := "select * from users"
	where, params := BuildFilter(filter)
	if len(where) > 0 {
		query = query + " where " + where
	}
	return query, params
}
func BuildFilter(filter *model.UserFilter) (string, []interface{}) {
	buildParam := s.BuildDollarParam
	var where []string
	var params []interface{}
	i := 1
	if len(filter.Id) > 0 {
		params = append(params, filter.Id)
		where = append(where, fmt.Sprintf(`id = %s`, buildParam(i)))
		i++
	}
	if filter.DateOfBirth != nil {
		if filter.DateOfBirth.Min != nil {
			params = append(params, filter.DateOfBirth.Min)
			where = append(where, fmt.Sprintf(`date_of_birth >= %s`, buildParam(i)))
			i++
		}
		if filter.DateOfBirth.Max != nil {
			params = append(params, filter.DateOfBirth.Max)
			where = append(where, fmt.Sprintf(`date_of_birth <= %s`, buildParam(i)))
			i++
		}
	}
	if len(filter.Username) > 0 {
		q := filter.Username + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`username like %s`, buildParam(i)))
		i++
	}
	if len(filter.Email) > 0 {
		q := filter.Email + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`email like %s`, buildParam(i)))
		i++
	}
	if len(filter.Phone) > 0 {
		q := "%" + filter.Phone + "%"
		params = append(params, q)
		where = append(where, fmt.Sprintf(`phone like %s`, buildParam(i)))
		i++
	}
	if len(where) > 0 {
		return strings.Join(where, " and "), params
	}
	return "", params
}
