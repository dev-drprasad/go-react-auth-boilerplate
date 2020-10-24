package main

import (
	"context"
	"flag"
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

func main() {
	var wait time.Duration
	var dbHost, dbPort, dbUser, dbPassword, dbName, port string
	flag.StringVar(&dbHost, "dbHost", "localhost", "")
	flag.StringVar(&dbName, "dbName", "repoboost", "")
	flag.StringVar(&dbUser, "dbUser", "repoboost", "")
	flag.StringVar(&dbPassword, "password", "repoboost", "")
	flag.StringVar(&dbPort, "dbPort", "5432", "")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

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

	if port = os.Getenv("PORT"); port == "" {
		port = "8000"
	}

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
