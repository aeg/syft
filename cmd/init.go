package cmd

import (
	"fmt"
	"os"

	"github.com/anchore/imgbom/imgbom"
	"github.com/anchore/imgbom/internal/config"
	"github.com/anchore/imgbom/internal/format"
	"github.com/anchore/imgbom/internal/log"
	"github.com/anchore/imgbom/internal/logger"
	"github.com/anchore/stereoscope"
	"github.com/spf13/viper"
	"github.com/wagoodman/go-partybus"
	"gopkg.in/yaml.v2"
)

var appConfig *config.Application
var eventBus *partybus.Bus
var eventSubscription *partybus.Subscription

func initAppConfig() {
	cfg, err := config.LoadConfigFromFile(viper.GetViper(), &cliOpts)
	if err != nil {
		fmt.Printf("failed to load application config: \n\t%+v\n", err)
		os.Exit(1)
	}
	appConfig = cfg
}

func initLogging() {
	config := logger.LogConfig{
		EnableConsole: (appConfig.Log.FileLocation == "" || appConfig.CliOptions.Verbosity > 0) && !appConfig.Quiet,
		EnableFile:    appConfig.Log.FileLocation != "",
		Level:         appConfig.Log.LevelOpt,
		Structured:    appConfig.Log.Structured,
		FileLocation:  appConfig.Log.FileLocation,
	}

	logWrapper := logger.NewZapLogger(config)
	imgbom.SetLogger(logWrapper)
	stereoscope.SetLogger(logWrapper)
}

func logAppConfig() {
	appCfgStr, err := yaml.Marshal(&appConfig)

	if err != nil {
		log.Debugf("Could not display application config: %+v", err)
	} else {
		log.Debugf("Application config:\n%+v", format.Magenta.Format(string(appCfgStr)))
	}
}

func initEventBus() {
	eventBus = partybus.NewBus()
	eventSubscription = eventBus.Subscribe()

	stereoscope.SetBus(eventBus)
	imgbom.SetBus(eventBus)
}
