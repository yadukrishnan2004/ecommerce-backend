package infrastructure

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectPostgres() (*gorm.DB,error){
	dsn:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_NAME"),
	os.Getenv("DB_PORT"),)
	DB,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err !=nil{
		return nil,errors.New("faile connecting data base")
	}
	return DB,nil
}