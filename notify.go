package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"net/smtp"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/http"
	"github.com/substitutes/substitutes-notify/mail"
)

var (
	verbose      = kingpin.Flag("verbose", "Enable verbose output").Short('v').Bool()
	daemon       = kingpin.Flag("daemon", "Run notify as a daemon").Short('d').Default("true").Bool()
	debugTrigger = kingpin.Flag("debugtrigger", "Always assume diff").Bool()
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

	viper.WatchConfig()

	if viper.ReadInConfig() != nil {
		log.Fatal("Failed to initialize configuration, does the config exist?")
	}

	log.Debug("Initialized application")

	viper.SetDefault("interval", "1h")
	viper.SetDefault("smtp_port", 25)

	viper.SetEnvPrefix("notify")
	viper.AutomaticEnv()

	// Connect to SMTP server

	auth := smtp.PlainAuth("", viper.GetString("smtp_username"), viper.GetString("smtp_password"), viper.GetString("smtp_host"))

	// Read in configuration file
	bytes, err := ioutil.ReadFile("users.json")

	if err != nil {
		log.Fatal("Failed to read users file: ", err)
	}

	var users Users
	if json.Unmarshal(bytes, &users) != nil {
		log.Fatal("Failed to marshal users json.")
	}

	ticker := time.NewTicker(viper.GetDuration("interval"))

	var classes []Data
	for range ticker.C {
		// TODO: Optimize scraping - maybe cached responses for classes already done in cycle?
		for _, u := range users {
			log.Debug("Fetching class " + u.Class)
			resp, err := http.Get(viper.GetString("api_url") + "/api/c/" + u.Class)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			var class Data
			if json.Unmarshal(bytes, &class) != nil {
				log.Fatal("Failed to marshal response!")
			}
			same := false
			for _, c := range classes {
				if c.Meta.Date == class.Meta.Date && c.Meta.Class == class.Meta.Class {
					same = true
				}
			}
			if !same || *debugTrigger {
				// Push notification to user
				update := mail.NewUpdate(viper.GetString("api_url")+"/c/"+u.Class, u.Class, u.Name, class.Meta.Date)
				updateMail := mail.New([]string{u.Email}, update, auth)
				if updateMail.Parse("mail/templates/update.html") != nil {
					log.Fatal("Failed to parse template: ", err)
				}
				log.Infof("Sent mail to %s (class %s (%s) updated [%s])", u.Name, u.Email, u.Class, class.Meta.Date)
				// updateMail.Send()
			}
			// Check if data is not overflowing
			classes = append(classes, class)
		}
	}
}
