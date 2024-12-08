package dbcontext

import (
	"allopopot-interconnect-service/config"
	"allopopot-interconnect-service/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	log.Println("Database Connection Initializing")
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		log.Panicln("Database Connection Failed")
	}
	DB = db
	log.Println("Database Connection Successful")
}

func Migrate() {
	log.Println("Database Migration Initializing")
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Panicln("Database Migration Failed")
	}
	log.Println("Database Migration Successful")

}
