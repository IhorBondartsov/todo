package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"to-do/app"
	"to-do/delivery"
	"to-do/repository"

	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var Version = "empty"

const (
	defaultAppName         = "todo-service"
	defaultDBDriver        = "postgres"
	defaultPort            = 8080
	defaultHost            = "0.0.0.0"
	defaultLogLevel        = "info"
	defaultShutdownTimeout = 10 * time.Second
)

type AppConfig struct {
	DB   repository.StorageConfig
	HTTP delivery.HTTPConfig

	AppName  string
	LogLevel string
}

func (cfg *AppConfig) Validate() error {
	errs := []error{}
	if err := cfg.DB.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := cfg.HTTP.Validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Errorf("%v", errs)
	}
	return nil
}

func parseAppCfg() (*AppConfig, error) {
	var config AppConfig

	flagset := flag.NewFlagSetWithEnvPrefix(defaultAppName, "", flag.ContinueOnError)

	// Global
	flagset.StringVar(&config.AppName, "app-name", defaultAppName, "Service name.")
	flagset.StringVar(&config.LogLevel, "log-level", defaultLogLevel, "Log level (debug, info, warn, error)")
	// HTTPConfig
	flagset.StringVar(&config.HTTP.Host, "host", defaultHost, "Host part of listening address.")
	flagset.IntVar(&config.HTTP.Port, "port", defaultPort, "Listening port.")
	flagset.DurationVar(&config.HTTP.ShutdownTimeout, "shutdown-timeout", defaultShutdownTimeout, "Shutdown timeout for http service.")
	//  DB
	flagset.StringVar(&config.DB.Driver, "db-driver", defaultDBDriver, "Data service driver.")
	flagset.StringVar(&config.DB.DSN, "db-dsn", "", "Data service data source name.")

	logrus.WithField("osenviron", os.Environ()).Info("env")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, errors.Wrap(err, "parsing flags")
	}

	// Validate the config.
	if err := config.Validate(); err != nil {
		return nil, errors.Errorf("invalid flag(s): %s", err)
	}
	return &config, nil
}

func initLogging(out io.Writer, level string) error {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("convert level: %s", err)
	}
	logger := logrus.Logger{
		Out:       out,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrusLevel,
	}

	log.SetOutput(logger.Writer())
	return nil
}

func main() {
	ctx := context.Background()
	log.Printf("running todo app, current version: %s \n", Version)

	cfg, err := parseAppCfg()
	if err != nil {
		logrus.Fatal(err)
	}

	err = initLogging(os.Stdout, cfg.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := repository.NewDBClient(ctx, cfg.DB)
	if err != nil {
		logrus.Fatal(err)
	}

	service, err := app.NewToDoService(db)
	if err != nil {
		logrus.Fatal(err)
	}

	httpService := delivery.NewHTTPService(cfg.HTTP, service)
	httpService.Run()
}
