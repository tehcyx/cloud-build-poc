package repository // import "github.com/tehcyx/cloud-build-poc/repository"

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// https://github.com/DATA-DOG/go-sqlmock
// https://github.com/jirfag/go-queryset/blob/master/queryset/queryset_test.go

//NewTestDB initializes the test-database
func NewTestDB() (sqlmock.Sqlmock, *DB) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("can't create sqlmock: %s", err)
	}
	db, err := gorm.Open("mysql", mockDb)
	if err != nil {
		log.Panic(fmt.Sprintf("something went wrong %s", err.Error()))
	}
	db.LogMode(true)
	return mock, &DB{db}
}
