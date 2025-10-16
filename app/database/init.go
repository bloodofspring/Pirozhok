package database

import (
	"errors"
	"main/database/models"
	"os"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var (
	db   *pg.DB
	once sync.Once
)

// GetDB returns a singleton instance of the database connection
func GetDB() *pg.DB {
	once.Do(func() {
		db = pg.Connect(&pg.Options{
			Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
			PoolSize: 20, // Устанавливаем разумный размер пула
		})
	})
	return db
}

func InitDb() error {
	db := GetDB()

	models := []interface{}{
		&models.Users{},
		&models.Groups{},
		&models.GroyupParticipants{},
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false, // Временные таблицы
			IfNotExists: true,
		})
		if err != nil {
			return errors.New("error creating table: " + err.Error())
		}
	}

	return nil
}
