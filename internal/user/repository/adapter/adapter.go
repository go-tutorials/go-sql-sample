package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	s "github.com/core-go/sql"

	"go-service/internal/user/model"
)

func NewUserAdapter(db *sql.DB, buildQuery func(*model.UserFilter) (string, []interface{})) (*UserAdapter, error) {
	userType := reflect.TypeOf(model.User{})
	fieldsIndex, _, jsonColumnMap, keys, _, _, buildParam, _, err := s.Init(userType, db)
	if err != nil {
		return nil, err
	}
	return &UserAdapter{DB: db, Map: fieldsIndex, Keys: keys, JsonColumnMap: jsonColumnMap, BuildParam: buildParam, BuildQuery: buildQuery}, nil
}

type UserAdapter struct {
	DB            *sql.DB
	Map           map[string]int
	Keys          []string
	JsonColumnMap map[string]string
	BuildParam    func(int) string
	BuildQuery    func(*model.UserFilter) (string, []interface{})
}

func (r *UserAdapter) All(ctx context.Context) ([]model.User, error) {
	query := `
		select
			id, 
			username,
			email,
			phone,
			date_of_birth
		from users`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.DateOfBirth)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (r *UserAdapter) Load(ctx context.Context, id string) (*model.User, error) {
	query := `
		select
			id, 
			username,
			email,
			phone,
			date_of_birth
		from users where id = $1`
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.DateOfBirth)
		return &user, nil
	}
	return nil, nil
}

func (r *UserAdapter) Create(ctx context.Context, user *model.User) (int64, error) {
	query := `
		insert into users (
			id,
			username,
			email,
			phone,
			date_of_birth)
		values (
			$1,
			$2,
			$3, 
			$4,
			$5)`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return -1, nil
	}
	res, err := stmt.ExecContext(ctx,
		user.Id,
		user.Username,
		user.Email,
		user.Phone,
		user.DateOfBirth)
	if err != nil {
		if strings.Index(err.Error(), "duplicate key") >= 0 {
			return -1, nil
		}
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Update(ctx context.Context, user *model.User) (int64, error) {
	query := `
		update users 
		set
			username = $1,
			email = $2,
			phone = $3,
			date_of_birth = $4
		where id = $5`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return -1, nil
	}
	res, err := stmt.ExecContext(ctx,
		user.Username,
		user.Email,
		user.Phone,
		user.DateOfBirth,
		user.Id)
	if err != nil {
		return -1, err
	}
	count, err := res.RowsAffected()
	if count == 0 {
		return count, nil
	}
	return count, err
}

func (r *UserAdapter) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	colMap := s.JSONToColumns(user, r.JsonColumnMap)
	query, args := s.BuildToPatch("users", colMap, r.Keys, r.BuildParam)
	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Delete(ctx context.Context, id string) (int64, error) {
	query := "delete from users where id = $1"
	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Search(ctx context.Context, filter *model.UserFilter, limit int64, offset int64) ([]model.User, int64, error) {
	var users []model.User
	if limit <= 0 {
		return users, 0, nil
	}
	query, params := r.BuildQuery(filter)
	pagingQuery := s.BuildPagingQuery(query, limit, offset)
	countQuery := s.BuildCountQuery(query)

	row := r.DB.QueryRowContext(ctx, countQuery, params...)
	if row.Err() != nil {
		return users, 0, row.Err()
	}
	var total int64
	err := row.Scan(&total)
	if err != nil || total == 0 {
		return users, total, err
	}

	err = s.Query(ctx, r.DB, r.Map, &users, pagingQuery, params...)
	return users, total, err
}

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
