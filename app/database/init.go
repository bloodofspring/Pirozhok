package database

import (
	"errors"
	"log"
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
	log.Println("Initializing database connection...")
	db := GetDB()

	// Проверяем подключение к базе данных
	_, err := db.Exec("SELECT 1")
	if err != nil {
		return errors.New("failed to connect to database: " + err.Error())
	}
	log.Println("Database connection established successfully")

	models := []interface{}{
		&models.GroupParticipants{},
		&models.Users{},
		&models.Groups{},
	}

	for _, model := range models {
		log.Printf("Creating table for model: %T", model)
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false, // Временные таблицы
			IfNotExists: true,
		})
		if err != nil {
			return errors.New("error creating table: " + err.Error())
		}
		log.Printf("Table created successfully for model: %T", model)
	}

	log.Println("Database initialization completed successfully")
	return nil
}
