package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// set http handlers
	mux := router.NewRouter(todoDB)

	// TODO: ここから実装を行う
	logChain := alice.New(middleware.GetOS, middleware.GetAccessLog)
	mux.Handle("/healthz", logChain.Then(handler.NewHealthzHandler()))
	hTODO := handler.NewTODOHandler(service.NewTODOService(todoDB))
	mux.Handle("/todos", logChain.Append(middleware.BasicAuth).Then(hTODO))
	hPanic := handler.NewPanicHandler()
	mux.Handle("/do-panic", logChain.Append(middleware.Recovery).Then(hPanic))
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	go func() {
		<-ctx.Done()

		log.Println("Graceful Shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}
