package user

import (
	"context"
	"database/sql"
	"fmt"
	. "go-service/internal/model"
	"strings"
)

type UserRepository interface {
	Load(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (int64, error)
	Update(ctx context.Context, user *User, id string) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}, id string) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{DB: db}
}

type userRepo struct {
	DB *sql.DB
}

func (s *userRepo) Load(ctx context.Context, id string) (*User, error) {
	query := "select * from usertest where id = $1 limit 1"
	rows, err := s.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	var res []User
	user := User{}
	for rows.Next() {
		err := rows.Scan(&user.Username, &user.Id, &user.Email, &user.Phone)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}

	return &res[0], nil
}

func (s *userRepo) Create(ctx context.Context, user *User) (int64, error) {
	query := fmt.Sprintf("INSERT INTO usertest(id, username, email, phone) VALUES ('%s', '%s', '%s', '%s')",
		user.Id,
		user.Username,
		user.Email,
		user.Phone)

	result, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *userRepo) Update(ctx context.Context, user *User, id string) (int64, error) {
	query := fmt.Sprintf("UPDATE usertest SET username='%s', email='%s', phone='%s' WHERE id='%s'",
		user.Username,
		user.Email,
		user.Phone,
		id)

	result, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *userRepo) Patch(ctx context.Context, user map[string]interface{}, id string) (int64, error) {
	updateClause := "UPDATE usertest SET"
	whereClause := fmt.Sprintf("WHERE id='%s'", id)

	setClause := make([]string, 0)
	if user["username"] != nil {
		msg := fmt.Sprintf("username='%s'", fmt.Sprint(user["username"]))
		setClause = append(setClause, msg)
	}
	if user["email"] != nil {
		msg := fmt.Sprintf("email='%s'", fmt.Sprint(user["email"]))
		setClause = append(setClause, msg)
	}
	if user["phone"] != nil {
		msg := fmt.Sprintf("phone='%s'", fmt.Sprint(user["phone"]))
		setClause = append(setClause, msg)
	}

	setClauseRes := strings.Join(setClause, ",")
	querySlice := []string{updateClause, setClauseRes, whereClause}
	query := strings.Join(querySlice, " ")

	result, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *userRepo) Delete(ctx context.Context, id string) (int64, error) {
	query := "delete from usertest where id = $1"

	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}