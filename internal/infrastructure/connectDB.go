package infrastructure

import (
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectPostgres(dsn string) (*gorm.DB){
	if dsn == "" {
		log.Fatal("no dsn found")
		return  nil
	}
	Dsn:=dsn
	DB,err:=gorm.Open(postgres.Open(Dsn),&gorm.Config{})
	if err !=nil{
		log.Fatal("faile to connect with the database")
		return nil
	}
	return DB
}