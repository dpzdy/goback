package utils

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func init() {
	dsn := "server=119.12.171.133;user id=topTenForeignNews;password=1503@cuc;port=1433;database=topTenForeignNews;encrypt=disable"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(); err != nil {
		log.Fatal(err)
	}
}
