package main

import (
	"context"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/Sxtanna/chromatic_curator/internal/app"
	"github.com/Sxtanna/chromatic_curator/internal/app/discord"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	"github.com/Sxtanna/chromatic_curator/internal/system/backend"
	"github.com/Sxtanna/chromatic_curator/internal/system/logging"
	"github.com/oklog/run"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"strings"
	"syscall"
)

type curatorConfiguration struct {
	Bot   *discord.BotConfiguration
	Log   *logging.Config
	Redis *backend.Config
}

func (c *curatorConfiguration) Validate() error {

	if err := common.OptValidate(c.Bot); err != nil {
		return err
	}

	if err := common.OptValidate(c.Log); err != nil {
		return err
	}

	return nil
}

func initializeConfiguration(v *viper.Viper, f *pflag.FlagSet) {
	v.AddConfigPath(".")

	v.SetConfigName(".env")
	v.SetConfigType("env")

	v.AllowEmptyEnv(true)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	_ = v.BindEnv("bot.token", "BOT_TOKEN")
	_ = v.BindEnv("bot.admins", "BOT_ADMINS")
	_ = v.BindEnv("redis.host", "REDIS_HOST")
	_ = v.BindEnv("redis.port", "REDIS_PORT")
}

func readPFlags(f *pflag.FlagSet) error {
	_ = f.Parse(os.Args[1:])

	if v, _ := f.GetBool("version"); v {
		os.Exit(0)
	}

	return nil
}

func readConfig(v *viper.Viper) error {
	var err error

	if err = v.ReadInConfig(); err != nil && !errors.Is(err, viper.ConfigFileNotFoundError{}) {
		return errors.Wrap(err, "failed to read config")
	}

	conf := curatorConfiguration{
		Bot:   &discord.BotConfiguration{},
		Log:   &logging.Config{},
		Redis: &backend.Config{},
	}

	if err = v.Unmarshal(&conf); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	} else {
		config = &conf
	}

	if err = common.OptProcess(config); err != nil {
		return errors.Wrap(err, "failed to process config")
	}

	logger, err = logging.New(*conf.Log)
	if err != nil {
		return errors.Wrap(err, "failed to create logger from config")
	}

	if err = common.OptValidate(config); err != nil {
		return errors.Wrap(err, "failed to validate config")
	}

	return nil
}

var (
	group run.Group
	abort chan struct{}

	config common.Configuration
	logger *slog.Logger
	handle emperror.ErrorHandler
)

func main() {
	var err error

	v, f := viper.New(), pflag.NewFlagSet("Chromatic Curator", pflag.ExitOnError)

	initializeConfiguration(v, f)

	f.Bool("version", false, "Show version")

	err = readPFlags(f)
	emperror.Panic(err)

	err = readConfig(v)
	emperror.Panic(err)

	logging.SetStandardLogger(logger)

	handle = logging.NewSlogHandler(logger)
	defer emperror.HandleRecover(handle)

	initializeRunGroup()

	err = initializeServices()
	emperror.Panic(err)

	logger.Info("application starting...")

	if err = group.Run(); !errors.As(err, &run.SignalError{}) {
		handle.Handle(err)
	}

	logger.Info("application closing...")
}

func initializeServices() error {
	serviceActors := app.InitializeApp(abort, logger, handle, config)
	for _, actor := range serviceActors {
		group.Add(actor.Execute, actor.Interrupt)
	}

	return nil
}

func initializeRunGroup() {
	abort = make(chan struct{})

	// signal handler
	execute, interrupt := run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	group.Add(execute, func(err error) {
		interrupt(err)

		logger.Debug("received interrupt, sending abort signal...")
		close(abort)
		logger.Debug("abort signal sent")
	})
}
