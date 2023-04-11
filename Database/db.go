package database

import (
	models "cricHeros/Models"
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	fmt.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.New("error in connecting to database")
	}
	query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`
	db.Exec(query)
	DB = db
	err = db.AutoMigrate(&models.Player{}, &models.Career{}, &models.Match{}, &models.Team{}, &models.Credential{}, &models.Balls{}, &models.ScoreCard{}, &models.MatchRecord{}, &models.TeamList{}, &models.Inning{}, &models.Toss{}, &models.Blacklist{})
	if err != nil {
		return errors.New("error in creating the tables")
	}
	fmt.Println("Succesfully connected to database...")
	return nil
}
