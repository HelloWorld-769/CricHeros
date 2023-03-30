package database

import (
	models "cricHeros/Models"
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
		fmt.Println("Error in connecting to database: ", err)
		return err
	}
	DB = db
	err = db.AutoMigrate(&models.Player{}, &models.Career{}, &models.Match{}, &models.Team{}, &models.Credential{}, &models.Balls{}, &models.ScoreCard{}, &models.MatchRecord{})
	if err != nil {
		fmt.Println("Error in creating the tables..")
		return err
	}
	fmt.Println("Succesfully connected to database...")
	return nil
}
