package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/skandyla/s3-uploader/internal/config"
	"github.com/skandyla/s3-uploader/internal/repository/psql"
	"github.com/skandyla/s3-uploader/internal/service"
	"github.com/skandyla/s3-uploader/internal/transport"
	"github.com/skandyla/s3-uploader/internal/version"
	"github.com/skandyla/s3-uploader/pkg/hash"
	"github.com/skandyla/s3-uploader/pkg/postgresql"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	initLogger(config.LogLevel, config.JsonLogOutput)
	log.Infof("starting app: %s, buildVersion: %s, buildBranch: %s, buildTime: %s, goVersion: %s", "microtester", version.BuildVersion, version.BuildBranch, version.BuildTime, version.GoVersion)

	log.Info("logLevel: ", config.LogLevel)
	log.Info("PostgresDSN: ", config.PostgresDSN)

	dbc, err := postgresql.NewConnection(config.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := dbc.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
		log.Println("clossing database connection")
	}()

	// init deps
	hasher := hash.NewSHA1Hasher("salt")

	healthRepoPg := psql.NewHealthRepository(dbc)
	healthService := service.NewHealth(healthRepoPg)

	tokensRepo := psql.NewTokens(dbc)
	usersRepo := psql.NewUsers(dbc)
	usersService := service.NewUsers(usersRepo, tokensRepo, hasher, []byte("sample secret"), config.Auth.TokenTTL)

	handler := transport.NewHandler(healthService, usersService)

	server := http.Server{
		Addr:           config.ListenAddress,
		Handler:        handler.InitRouter(),
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server started, listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("server error: %v", err)
		return

	case sig := <-quit:
		log.Info("shutdown started, signal: ", sig)
		defer log.Info("shutdown complete, signal: ", sig)

		ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			log.Printf("could not stop server gracefully: %v", err)
			return
		}
	}
}

//------------------------------
func initLogger(logLevel string, json bool) {
	if json {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetOutput(os.Stderr)

	switch strings.ToLower(logLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}
