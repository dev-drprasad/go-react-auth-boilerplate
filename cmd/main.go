package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	authrest "repoboost/internal/auth/rest"
	userrest "repoboost/internal/user/rest"

	"github.com/gorilla/mux"
	pgxpool "github.com/jackc/pgx/v4/pgxpool"
)

func getEnvDefault(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {

	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	dbHost := getEnvDefault("DB_HOST", "localhost")
	dbPort := getEnvDefault("DB_PORT", "5432")
	dbUser := getEnvDefault("DB_USER", "repoboost")
	dbPassword := getEnvDefault("DB_PASSWORD", "repoboost")
	dbName := getEnvDefault("DB_NAME", "repoboost")
	port := getEnvDefault("PORT", "8000")

	connstr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, connstr)
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %v", err))
	}
	defer db.Close()

	userhandler := userrest.New(db)
	authhandler := authrest.New(db)

	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/auth/login", authhandler.Login).Methods(http.MethodPost)

	v1.HandleFunc("/users", userhandler.GetUsers).Methods(http.MethodGet)
	v1.HandleFunc("/users", userhandler.CreateUser).Methods(http.MethodPost)
	v1.HandleFunc("/users", userhandler.GetUser).Methods(http.MethodGet)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("server started at %s\n", port)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
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
