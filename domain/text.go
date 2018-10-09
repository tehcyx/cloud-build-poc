package domain

// Text represents a single text entity
type Text struct {
	Model
	Text  string
	Title string

	Store TextStore `gorm:"-" json:"-"` // Interface
}

// NewText create a new text object
func NewText(text, title string) *Text {
	t := &Text{}
	if text == "" || title == "" {
		return nil
	}
	t.Text = text
	t.Title = title
	return t
}

//TextStore interface
type TextStore interface {
	CreateText(*Text)
	GetText(uint) *Text
	GetTexts() []*Text
	UpdateText(*Text) *Text
	DeleteText(uint) error
}

//NewTextStore returns an instance with the db/store dependency met
func NewTextStore(ts TextStore) *Text {
	return &Text{
		Store: ts,
	}
}
