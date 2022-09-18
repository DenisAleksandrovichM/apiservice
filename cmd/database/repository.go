package main

import (
	"fmt"
	configPkg "github.com/DenisAleksandrovichM/apiservice/internal/database/config"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func createRepository() (*sqlx.DB, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configPkg.DBHost, configPkg.DBPort, configPkg.DBUser, configPkg.DBPassword, configPkg.DBName)

	db, err := sqlx.Open("pgx", psqlConn)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to database")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping database error")
	}
	return db, nil
}
