package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

var (
	verbose = kingpin.Flag("verbose", "Enable verbose output").Short('v').Bool()
	daemon = kingpin.Flag("daemon", "Run notify as a daemon").Short('d').Default("true").Bool()
)

func main() {
	kingpin.Parse()

	log.SetLevel(log.WarnLevel)
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/substitutes/")
	viper.AddConfigPath("/etc/substitutes/")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if viper.ReadInConfig() != nil {
		log.Fatal("Failed to initialize configuration, does the config exist?")
	}

	log.Debug("Initialized application")

}