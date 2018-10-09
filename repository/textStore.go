package repository

import "github.com/tehcyx/cloud-build-poc/domain"

// CreateText creates a Text inside the repository
func (db *DB) CreateText(c *domain.Text) {
	db.Create(c)
}

// GetText finds the Text by it's unique id
func (db *DB) GetText(id uint) []*domain.Text {
	var textData []*domain.Text
	db.First(&textData, id)
	return textData
}

// GetTexts finds all Texts inside the repository
func (db *DB) GetTexts() []*domain.Text {
	var textData []*domain.Text
	db.Order("created_at desc").Find(&textData)
	return textData
}

// UpdateText updates specific text inside repository
func (db *DB) UpdateText(c *domain.Text) *domain.Text {
	db.Save(c)
	return c
}

// DeleteText deletes specific text by ID inside repository
func (db *DB) DeleteText(id uint) error {
	err := db.Where("id = ?", id).Delete(domain.Text{}).Error
	return err
}
