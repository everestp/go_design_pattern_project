package models

import "database/sql"


type Repository interface{
	AllDogBreeds() ([]*DogBreed , error)
}

type mysqlRepository struct{
	DB *sql.DB
}


func newMysqlRepository(conn *sql.DB) Repository{
	return  &mysqlRepository{
		DB: conn,
	}
}