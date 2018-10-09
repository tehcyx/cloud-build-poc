package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/tehcyx/cloud-build-poc/src/domain"
	"github.com/tehcyx/cloud-build-poc/src/repository"
	"github.com/tehcyx/cloud-build-poc/src/text"
	textBoundary "github.com/tehcyx/cloud-build-poc/src/text/boundary"
)

var db *repository.DB
var logger *log.Logger
var router *mux.Router
var applicationName = "cloud-builder-poc"

// InitDB initializes the postgresdb connection with a stupid backoff retry
func InitDB() *repository.DB {
	// GET ENV VARS
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DBNAME")
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPass == "" {
		dbPass = "secret"
	}
	if dbName == "" {
		dbName = "web"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	dbHandle, err := repository.NewPostgresDB(psqlInfo, 5, logger)
	if err != nil {
		logger.Fatalf("something went wrong connecting to the database: %s\n", err.Error())
	}
	// Disable table name's pluralization globally
	dbHandle.SingularTable(true) // if set this to true, `User`'s default table name will be `user`, table name setted with `TableName` won't be affected

	text.InitShared(dbHandle, logger)
	domain.InitShared(logger)

	return dbHandle
}

// InitLogger initializes the logger format to use throughout the application
func InitLogger() *log.Logger {
	logHandle := log.New(os.Stdout, fmt.Sprintf("%s ", applicationName), log.LstdFlags|log.Lshortfile)
	return logHandle
}

// InitRouter initializes the router and registers the different subrouters to it which have their own initializers
func InitRouter() *mux.Router {
	routeHandler := mux.NewRouter().StrictSlash(true)
	routeHandler.NotFoundHandler = domain.AppHandler(domain.NotFoundHandler)

	textBoundary.RegisterTextRouter(routeHandler)

	return routeHandler
}

// InitHandler glues together all initializers and exposes the server on the defined port. Handles a graceful shutdown of the program for SIGKILL, SIGQUIT and SIGTERM
func InitHandler() {
	router = InitRouter()
	logger = InitLogger()
	db = InitDB()

	// if local: firewall of mac won't ask for permission when 127.0.0.1/localhost is specified
	host := "0.0.0.0"
	port := ":8080"

	if os.Getenv("DEPLOY_ENV") == "" {
		logger.Printf("Running dev environment\n")
		host = "localhost"
	}

	addr := fmt.Sprintf("%s%s", host, port)

	srv := &http.Server{
		Handler: domain.CORS().Handler(router),
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		// from: https://blog.cloudflare.com/exposing-go-on-the-internet/
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m
	var wait time.Duration
	wait = time.Second * 15

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logger.Println(fmt.Sprintf("Listening on %s", port))
		if err := srv.ListenAndServe(); err != nil {
			logger.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
