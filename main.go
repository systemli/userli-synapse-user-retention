package main

import (
	"errors"
	"os"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"
)

func init() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "text" {
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithError(err).Fatal("Failed to parse log level")
	}
	log.SetLevel(level)
}

func main() {
	if err := checkEnvironment(); err != nil {
		log.WithError(err).Fatal("Failed to check environment")
	}

	synapse := NewSynapseClient(os.Getenv("SYNAPSE_URL"), os.Getenv("SYNAPSE_TOKEN"))
	userli := NewUserliClient(os.Getenv("USERLI_URL"), os.Getenv("USERLI_DOMAIN"), os.Getenv("USERLI_TOKEN"))

	synapseUsers, err := synapse.FetchUsers()
	if err != nil {
		log.WithError(err).Fatal("Failed to fetch Synapse users")
	}
	log.WithField("users", len(synapseUsers)).Info("Fetched Synapse users")

	userliDeletedUsers, err := userli.FetchDeletedUsers()
	if err != nil {
		log.WithError(err).Fatal("Failed to fetch deleted Userli users")
	}
	log.WithField("users", len(userliDeletedUsers)).Info("Fetched deleted Userli users")

	for _, user := range synapseUsers {
		email := strings.ReplaceAll(user.Name[1:], ":", "@")
		if slices.Contains(userliDeletedUsers, email) {
			log.WithField("email", email).Info("Deleting user")
			// TODO: Delete user from Synapse
		} else {
			if user.LastSeen > 0 {
				log.WithFields(log.Fields{"email": email, "timestamp": user.LastSeen}).Info("Touching user")
				if err := userli.TouchUser(email, user.LastSeen); err != nil {
					log.WithError(err).WithField("email", email).Error("Failed to touch user")
					continue
				}
			}
		}
	}
}

func checkEnvironment() error {
	if os.Getenv("SYNAPSE_URL") == "" {
		return errors.New("SYNAPSE_URL is missing")
	}

	if os.Getenv("SYNAPSE_TOKEN") == "" {
		return errors.New("SYNAPSE_TOKEN is missing")
	}

	if os.Getenv("USERLI_URL") == "" {
		return errors.New("USERLI_URL is missing")
	}

	if os.Getenv("USERLI_DOMAIN") == "" {
		return errors.New("USERLI_DOMAIN is missing")
	}

	if os.Getenv("USERLI_TOKEN") == "" {
		return errors.New("USERLI_TOKEN is missing")
	}

	return nil
}
