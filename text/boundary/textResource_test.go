package boundary_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tehcyx/cloud-build-poc/domain"
	"github.com/tehcyx/cloud-build-poc/repository"
	"github.com/tehcyx/cloud-build-poc/text"
	textBoundary "github.com/tehcyx/cloud-build-poc/text/boundary"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var db *repository.DB
var logger *log.Logger
var router *mux.Router
var applicationName = "cloud-builder-poc-test"

var sqlMock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	initMockShared()
	initMockRouter()

	os.Exit(m.Run())
}

// func TestCORSEnabled(t *testing.T) {
// 	req, _ := http.NewRequest("OPTIONS", "/texts/v1/text", nil)

// 	resp := executeRequest(req)

// 	checkResponseCode(t, http.StatusOK, resp.Code)

// 	log.Printf("%v", resp.HeaderMap)

// 	// AllowedOrigins: []string{"http://localhost:8080"},
// 	// 	AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT"},
// 	// 	AllowedHeaders: []string{"Content-Type", "Content-Length", "Cache-Control", "X-Requested-With"},
// }

func TestTextNotFoundHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/texts/", nil)

	resp := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, resp.Code)

	var m map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &m)

	checkValues(t, "Not Found", m["message"].(string))
	checkValues(t, "error", m["type"].(string))
}

func TestTextV1RouterAccessible(t *testing.T) {
	req, _ := http.NewRequest("GET", "/texts/v1/text", nil)

	resp := executeRequest(req)

	checkResponseCode(t, http.StatusOK, resp.Code)
}

func TestTextV1GetTextByIDNumberAccessible(t *testing.T) {
	req, _ := http.NewRequest("GET", "/texts/v1/text/1", nil)

	resp := executeRequest(req)

	checkResponseCode(t, http.StatusOK, resp.Code)
}

func TestTextV1GetTextByIDStringBadRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/texts/v1/text/string", nil)

	resp := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, resp.Code)

	var m map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &m)

	checkValues(t, "something went wrong", m["message"].(string))
	checkValues(t, "error", m["type"].(string))
}

func TestTextV1GetTextByIDBadIntegerBadRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/texts/v1/text/-1", nil)

	resp := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, resp.Code)

	var m map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &m)

	checkValues(t, "something went wrong", m["message"].(string))
	checkValues(t, "error", m["type"].(string))
}

func checkValues(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected response code %s. Got %s\n", expected, actual)
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func initMockShared() {
	logger = log.New(os.Stdout, fmt.Sprintf("%s ", applicationName), log.LstdFlags|log.Lshortfile)

	// var mock sqlmock.Sqlmock
	sqlMock, db = repository.NewTestDB()

	text.InitShared(db, logger)
	domain.InitShared(logger)
}

func initMockRouter() {
	routeHandler := mux.NewRouter().StrictSlash(true)
	routeHandler.NotFoundHandler = domain.AppHandler(domain.NotFoundHandler)

	textBoundary.RegisterTextRouter(routeHandler)

	router = routeHandler
}
