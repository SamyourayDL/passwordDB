package main

import (
	"fmt"
	"log/slog"
	"os"
	"password-db/internals/config"
	"password-db/internals/lib/logger/slogpretty"
	"password-db/internals/storage/postgres"
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

	// secret, _ := crypto.Encrypt([]byte("lagoflash"), cfg.Key)
	// _, _ = secret, storage
	// fmt.Printf("%x\n", secret)

	resp, err := storage.GetPass("Alice", "")
	if err != nil {
		log.Error("failed to GetPass", "err", err.Error())
		os.Exit(1)
	}
	fmt.Println(resp)

	// err = storage.AddUser("Alice")
	// if err != nil {
	// 	log.Error("failed to AddUser", "err", err.Error())
	// 	os.Exit(1)
	// }

	// err = storage.AddPassword("Alice", "1211990", "sber", "banks")
	// if err != nil {
	// 	log.Error("failed to AddUser", "err", err.Error())
	// 	os.Exit(1)
	// }

	rowsDeleted, err := storage.Delete("Alice", "")
	if err != nil {
		log.Error("failed to delete user", "err", err.Error())
		os.Exit(1)
	}
	fmt.Println(rowsDeleted)

	log.Info("storage has been successfully inited")

	//router

	//listen and serve
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
