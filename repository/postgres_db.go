package repository // import "github.com/tehcyx/cloud-build-poc/repository"

import (
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

//NewPostgresDB initializes the database
func NewPostgresDB(psqlInfo string, maxReconnects int, logHandle *log.Logger) (*DB, error) {
	var db *gorm.DB
	var err error
	for try := 0; try < maxReconnects; try++ {
		db, err = gorm.Open("postgres", psqlInfo)
		if err != nil && try < maxReconnects {
			backoff := rand.Intn(5)
			logHandle.Printf("retrying db connection in %d seconds", backoff)
			time.Sleep(time.Second * time.Duration(backoff))
		}
	}
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
