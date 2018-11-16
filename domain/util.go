package domain // import "github.com/tehcyx/cloud-build-poc/domain"

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/cors"
)

var logger *log.Logger

// InitShared shared data initializer for control package
func InitShared(logHandle *log.Logger) {
	logger = logHandle
}

// appHandler is to be used in error handling
type AppHandler func(http.ResponseWriter, *http.Request) *AppError

type AppError struct {
	Err     error
	Message string
	Code    int
}

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *AppError
		logger.Println(e.Err)
		err := &map[string]string{
			"status":  strconv.Itoa(e.Code),
			"type":    "error",
			"message": e.Message,
		}
		RespondWithJSON(w, err, e.Code)
	}
}

// RespondWithJSON HandleFunc to create json responses from arbitrary payloads
func RespondWithJSON(w http.ResponseWriter, payload interface{}, code int) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondNoContent HandleFunc to create a no content response
func RespondNoContent(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write(nil)
}

// NotFoundHandler custom json not found error response
func NotFoundHandler(w http.ResponseWriter, r *http.Request) *AppError {
	return &AppError{Err: fmt.Errorf("Not Found"), Message: "Not Found", Code: http.StatusNotFound}
}

func MiddlewareHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}

type IndexHandler struct{}

func (ih *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/index.html")
}

func NoDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			logger.Println(r.URL.Path, http.StatusNotFound, "Not Found")
			err := &map[string]string{
				"status":  strconv.Itoa(http.StatusNotFound),
				"type":    "error",
				"message": "Not Found",
			}
			RespondWithJSON(w, err, http.StatusNotFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// CORS middleware adding cors headers to req
func CORS() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT"},
		AllowedHeaders: []string{"Content-Type", "Content-Length", "Cache-Control", "X-Requested-With"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	return c
}

// TimeTrack functions to measure execution time.
// usage: defer util.TimeTrack(time.Now(), "function")
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// RandomString returns a random string with the specified length
func RandomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
