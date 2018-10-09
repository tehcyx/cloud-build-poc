package text

import (
	"log"

	"github.com/tehcyx/cloud-build-poc/src/domain"
	"github.com/tehcyx/cloud-build-poc/src/repository"
	"github.com/tehcyx/cloud-build-poc/src/text/boundary"
	"github.com/tehcyx/cloud-build-poc/src/text/control"
)

var db *repository.DB
var logger *log.Logger

// InitShared initializes db and logger with the references passed and passes it to subpackages
func InitShared(dbHandle *repository.DB, logHandle *log.Logger) {
	db = dbHandle
	db.AutoMigrate(&domain.Text{})

	logger = logHandle

	control.InitShared(db, logger)
	boundary.InitShared(db, logger)
}

// IsInitialized returns wether everything was initialized successfully (true) or if it holds nil (false)
func IsInitialized() bool {
	return db != nil && logger != nil
}
