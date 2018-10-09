package text_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/tehcyx/cloud-build-poc/src/repository"
	"github.com/tehcyx/cloud-build-poc/src/text"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var sqlMock sqlmock.Sqlmock
var testDBHandle *repository.DB
var logHandle *log.Logger
var applicationName = "cloud-builder-poc-test"

func TestMain(m *testing.M) {
	logHandle = log.New(os.Stdout, fmt.Sprintf("%s ", applicationName), log.LstdFlags|log.Lshortfile)

	// var mock sqlmock.Sqlmock
	sqlMock, testDBHandle = repository.NewTestDB()
	text.InitShared(testDBHandle, logHandle)

	os.Exit(m.Run())
}

func TestInitializeSuccessful(t *testing.T) {
	if !text.IsInitialized() {
		t.Errorf("there were unfulfilled expectations: %s", errors.New("Either logger or db was not succesfully initialized"))
	}
}
