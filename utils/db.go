package utils

import (
	"goback/models"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
)

func InitDB() *gorm.DB {
	dsn := "server=119.12.171.133;user id=topTenForeignNews;password=1503@cuc;port=1433;database=topTenForeignNews;encrypt=disable"
	DB, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := DB.AutoMigrate(&models.Info{}); err != nil {
		log.Fatal(err)
	}
	return DB
}
