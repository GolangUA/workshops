package postgre

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Roma7-7-7/workshops/calendar/internal/models"
)

func (r *Repository) GetUser(name string) (*models.User, error) {
	var user models.User
	err := psql.Select("name", "password", "timezone").
		From("users").
		Where(sq.Eq{"name": name}).
		RunWith(r.db).
		QueryRow().
		Scan(&user.Name, &user.Password, &user.Timezone)
	if err != nil {
		return nil, fmt.Errorf("get user: %v", err)
	}
	return &user, nil
}

func (r *Repository) CreateUser(name string, password string, timezone string) (*models.User, error) {
	var user models.User
	err := psql.Insert("users").
		Columns("name", "password", "timezone").
		Values(name, password, timezone).
		Suffix("RETURNING name, password, timezone").
		RunWith(r.db).
		QueryRow().
		Scan(&user.Name, &user.Password, &user.Timezone)
	if err != nil {
		return nil, fmt.Errorf("create user: %v", err)
	}

	return &user, nil
}

func (r *Repository) UpdateUserTimezone(name string, timezone string) (*models.User, error) {
	var user models.User
	err := psql.Update("users").
		Set("timezone", timezone).
		Where(sq.Eq{"name": name}).
		Suffix("RETURNING name, password, timezone").
		RunWith(r.db).
		QueryRow().
		Scan(&user.Name, &user.Password, &user.Timezone)
	if err != nil {
		return nil, fmt.Errorf("update user timezone: %v", err)
	}

	return &user, nil
}

func (r *Repository) GetEventOwner(eventId string) (string, error) {
	query, args, err := psql.Select("username").
		From("user_event").
		Where(sq.Eq{"event_id": eventId}).
		ToSql()

	if err != nil {
		return "", fmt.Errorf("build event owner query: %v", err)
	}
	var username string
	if err = r.db.QueryRow(query, args...).Scan(&username); err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("scan event owner: %v", err)
	}

	return username, nil
}
