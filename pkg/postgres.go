package pkg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)


func NewPostgresDB(cfg DBconfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", 
				cfg.Host, cfg.Port, cfg.Username, cfg.DBname, cfg.Password, cfg.SSLMode))
	if err != nil {
		logrus.Fatalf("Error of open DB %s", err.Error())
		return nil, err
	}
	if err = db.Ping(); err != nil {
		logrus.Fatalf("Ping error: %s", err.Error())
		return nil, err
	}
	return db, nil
}

