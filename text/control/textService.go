package control

import (
	"log"

	"github.com/tehcyx/cloud-build-poc/domain"
	"github.com/tehcyx/cloud-build-poc/repository"
)

var db *repository.DB
var logger *log.Logger

// InitShared shared data initializer for text package
func InitShared(dbHandle *repository.DB, logHandle *log.Logger) {
	db = dbHandle
	logger = logHandle
}

// GetAllTexts returns list of all texts
func GetAllTexts() []*domain.Text {
	var texts = db.GetTexts()
	return texts
}

// GetText returns text received by id or nil
func GetText(textID uint) *domain.Text {
	var text = db.GetText(textID)

	if text == nil || len(text) == 0 {
		return nil
	}
	return text[0]
}
