package database

import (
	"../models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(){
	connection, err := gorm.Open(mysql.Open("username:password@/databasename"), &gorm.Config{})

	if err != nil{
		panic("Database Connection Failed")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
}

