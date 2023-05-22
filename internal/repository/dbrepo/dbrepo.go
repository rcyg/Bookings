package dbrepo

import (
	"Bookings/internal/config"
	"Bookings/internal/repository"
	"database/sql"
)

type postgreDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgreRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgreDBRepo{
		App: a,
		DB:  conn,
	}
}
