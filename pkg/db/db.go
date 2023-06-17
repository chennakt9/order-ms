package db

import (
	"fmt"
	"log"

	"github.com/chennakt9/order-ms/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Order{})
	
	fmt.Println("DB connected successfully")
	
	return Handler{db}
}