package domain // import "github.com/tehcyx/cloud-build-poc/domain"

import "time"

// Model represents basic fields for all entries
type Model struct {
	ID        uint       `gorm:"primary_key,auto_increment" json:"id" valid:"-"`
	CreatedAt time.Time  `json:"-" valid:"-"`
	UpdatedAt time.Time  `json:"-" valid:"-"`
	DeletedAt *time.Time `sql:"index" json:"-" valid:"-"`
}
