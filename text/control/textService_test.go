package control

import (
	"os"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/tehcyx/cloud-build-poc/repository"
)

var sqlMock sqlmock.Sqlmock
var testDBHandle *repository.DB

func TestMain(m *testing.M) {
	// var mock sqlmock.Sqlmock
	sqlMock, testDBHandle = repository.NewTestDB()
	InitShared(testDBHandle, nil)

	os.Exit(m.Run())
}

func TestGetAllTextsSuccess(t *testing.T) {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "text", "title"}).
		AddRow(1, time.Now(), time.Now(), nil, "hello", "world").
		AddRow(1, time.Now(), time.Now(), nil, "hello", "world")
	sqlMock.ExpectQuery("^SELECT (.+) FROM `texts` WHERE `texts`.`deleted_at` IS NULL ORDER BY created_at desc$").WillReturnRows(rows)

	texts := GetAllTexts()

	// we make sure that all expectations were met
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if len(texts) != 2 {
		t.Errorf("Missing rows")
	}
}

func TestGetAllTextsSuccessEmpty(t *testing.T) {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "text", "title"})
	sqlMock.ExpectQuery("^SELECT (.+) FROM `texts` WHERE `texts`.`deleted_at` IS NULL ORDER BY created_at desc$").WillReturnRows(rows)

	texts := GetAllTexts()

	// we make sure that all expectations were met
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if len(texts) != 0 {
		t.Errorf("Missing rows")
	}
}

func TestGetTextByIdSuccess(t *testing.T) {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "text", "title"}).
		AddRow(1, time.Now(), time.Now(), nil, "hello", "world")

	sqlMock.ExpectQuery("^SELECT (.+) FROM `texts`").WillReturnRows(rows)

	text := GetText(1)

	// we make sure that all expectations were met
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if text == nil {
		t.Errorf("Text can't be nil")
	}
}

func TestGetTextByIdSuccessEmpty(t *testing.T) {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "text", "title"})

	sqlMock.ExpectQuery("^SELECT (.+) FROM `texts`").WillReturnRows(rows)

	text := GetText(1)

	// we make sure that all expectations were met
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if text != nil {
		t.Errorf("Text should be nil")
	}
}
