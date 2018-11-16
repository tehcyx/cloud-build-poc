package repository // import "github.com/tehcyx/cloud-build-poc/repository"

import "github.com/jinzhu/gorm"

// DB pass around to have one single DB instance initialized on startup
type DB struct {
	*gorm.DB
}
