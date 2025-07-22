package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"password-db/internals/config"
	"password-db/internals/http-server/logger"
	"password-db/internals/http-server/password/deletepass"
	"password-db/internals/http-server/password/getpass"
	"password-db/internals/http-server/password/postpass"
	"password-db/internals/http-server/user/deleteuser"
	"password-db/internals/http-server/user/getuser"
	"password-db/internals/http-server/user/postuser"

	"password-db/internals/lib/logger/slogpretty"
	"password-db/internals/storage/postgres"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//parse config
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := SetupLogger(cfg.Env)
	log.Info("config has been successfully read")
	log.Info("logger has been successfully setuped")

	storage, err := postgres.New(&cfg.DB)
	if err != nil {
		log.Error("failed to init storage", "err", err.Error())
		os.Exit(1)
	}
	_ = storage

	log.Info("storage has been successfully inited")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// router.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	router.Route("/user", func(r chi.Router) {
		r.Post("/{user_name}", postuser.New(log, storage))
		r.Get("/{user_name}", getuser.New(log, storage))
		r.Delete("/{user_name}", deleteuser.New(log, storage))
	})

	router.Route("/password", func(r chi.Router) {
		r.Post("/{user_name}", postpass.New(log, storage))
		r.Get("/{user_name}", getpass.New(log, storage))
		r.Delete("/{user_name}", deletepass.New(log, storage))
	})

	srv := http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Error("failed to run server", "err", err.Error())
		os.Exit(1)
	}
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelError,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func testStorage(storage *postgres.Storage, log *slog.Logger) {
	// secret, _ := crypto.Encrypt([]byte("lagoflash"), cfg.Key)
	// _, _ = secret, storage
	// fmt.Printf("%x\n", secret)

	resp, err := storage.GetPass("Alice", "")
	if err != nil {
		log.Error("failed to GetPass", "err", err.Error())
		os.Exit(1)
	}
	fmt.Println(resp)

	err = storage.AddUser("Alice")
	if err != nil {
		log.Error("failed to AddUser", "err", err.Error())
		os.Exit(1)
	}

	err = storage.AddPassword("Alice", "1211990", "sber", "banks")
	if err != nil {
		log.Error("failed to AddUser", "err", err.Error())
		os.Exit(1)
	}

	rowsDeleted, err := storage.Delete("Alice", "")
	if err != nil {
		log.Error("failed to delete user", "err", err.Error())
		os.Exit(1)
	}
	fmt.Println(rowsDeleted)
}
