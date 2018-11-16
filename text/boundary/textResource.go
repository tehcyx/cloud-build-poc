package boundary // import "github.com/tehcyx/cloud-build-poc/text/boundary"

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tehcyx/cloud-build-poc/domain"
	"github.com/tehcyx/cloud-build-poc/repository"
	"github.com/tehcyx/cloud-build-poc/text/control"
	"github.com/urfave/negroni"
)

const (
	pathPrefix = "/texts"
)

var db *repository.DB
var logger *log.Logger

// InitShared shared data initializer for text package
func InitShared(dbHandle *repository.DB, logHandle *log.Logger) {
	db = dbHandle
	logger = logHandle
}

// RegisterTextRouter Pass a router and this function will automatically handle registration of the user resources to this router.
func RegisterTextRouter(r *mux.Router) {
	textRouter := textResourceRouter()
	r.PathPrefix(pathPrefix).Handler(negroni.New(
		negroni.HandlerFunc(domain.MiddlewareHandler),
		negroni.Wrap(textRouter),
	))
}

// textResourceRouter mux router for the `/texts` resource
// Make sure that the app config has the same PathPrefix as specified here.
func textResourceRouter() *mux.Router {
	textRouter := mux.NewRouter().PathPrefix(pathPrefix).Subrouter().StrictSlash(true)
	textRouter.NotFoundHandler = domain.AppHandler(domain.NotFoundHandler)

	textRouter.Handle("/v1/text", domain.CORS().Handler(domain.AppHandler(textV1RootHandler))).Methods(http.MethodGet, http.MethodOptions)
	textRouter.Handle("/v1/text/{textId}", domain.CORS().Handler(domain.AppHandler(textV1GetTextByIDHandler))).Methods(http.MethodGet, http.MethodOptions)
	return textRouter
}

func textV1RootHandler(w http.ResponseWriter, r *http.Request) *domain.AppError {
	defer domain.TimeTrack(time.Now(), "textV1RootHandler")

	texts := control.GetAllTexts()

	domain.RespondWithJSON(w, texts, http.StatusOK)

	return nil
}

func textV1GetTextByIDHandler(w http.ResponseWriter, r *http.Request) *domain.AppError {
	defer domain.TimeTrack(time.Now(), "textV1GetTextByIDHandler")

	vars := mux.Vars(r)

	tid, err := strconv.ParseUint(vars["textId"], 10, 64)
	if err != nil {
		return &domain.AppError{Err: errors.New("something went wrong"), Code: http.StatusBadRequest, Message: "something went wrong"}
	}
	textID := uint(tid)

	text := control.GetText(textID)

	domain.RespondWithJSON(w, text, http.StatusOK)

	return nil
}
