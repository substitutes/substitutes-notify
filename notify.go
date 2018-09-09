package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"net/smtp"
	"github.com/fronbasal/substitutes-notify/mail"
)

var (
	verbose = kingpin.Flag("verbose", "Enable verbose output").Short('v').Bool()
	daemon  = kingpin.Flag("daemon", "Run notify as a daemon").Short('d').Default("true").Bool()
)

var (
	auth smtp.Auth
)

func main() {
	kingpin.Parse()
	kingpin.CommandLine.Author("Daniel Malik <mail@fronbasal.de>")
	kingpin.CommandLine.Name = "Substitutes Notify"

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

	viper.SetDefault("interval", "1h")
	viper.SetDefault("smtp_port", 25)

	viper.SetEnvPrefix("notify")
	viper.AutomaticEnv()

	// Connect to SMTP server

	auth = smtp.PlainAuth("", viper.GetString("smtp_username"), viper.GetString("smtp_password"), viper.GetString("smtp_host"))

	req := mail.New([]string{"daniel.malik@steinbart-gym.eu"}, mail.NewUpdate("https://example.com", "11", "Daniel", "24.24.24"), auth)
	if err := req.Parse("mail/templates/update.html"); err != nil {
		log.Fatal("Failed to parse template: ", err)
		return
	}
	if err := req.Send(); err != nil {
		log.Fatal("Failed to send mail: ", err)
		return
	}
}
