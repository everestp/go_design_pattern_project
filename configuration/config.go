package configuration

import (
	"database/sql"
	"go-breeders/models"
	"sync"
)

type Application struct{
	Models *models.Models

}


var instance *Application

var once sync.Once
var db *sql.DB

func New(pool *sql.DB) *Application{
	db = pool
	return  GetIntance()
}

func GetIntance() *Application{
 once.Do(func(){
	instance = &Application{Models: models.New(db)}
 })

 return  instance
}
