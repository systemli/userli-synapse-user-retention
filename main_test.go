package main

import (
	"io"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
}

func (s *MainTestSuite) SetupTest() {
	log.SetOutput(io.Discard)
	defer gock.Off()
}

func (s *MainTestSuite) TestMain() {
	s.Run("Success", func() {
		os.Setenv("SYNAPSE_URL", "http://synapse")
		os.Setenv("SYNAPSE_TOKEN", "token")
		os.Setenv("USERLI_URL", "http://userli")
		os.Setenv("USERLI_DOMAIN", "example.com")
		os.Setenv("USERLI_TOKEN", "token")

		gock.New("http://synapse").
			Get("/_synapse/admin/v3/users").
			Reply(200).
			JSON(map[string]interface{}{
				"users": []map[string]interface{}{
					{
						"name":         "@test:example.com",
						"creation_ts":  123456,
						"last_seen_ts": 123456,
					},
				},
			})

		gock.New("http://userli").
			Get("/api/retention/example.com/users").
			Reply(200).
			JSON([]string{})

		gock.New("http://userli").
			Put("/api/retention/test@example.com/touch").
			Reply(200)

		main()

		s.True(gock.IsDone())
	})
}

func (s *MainTestSuite) TestCheckEnvironment() {
	s.Run("Success", func() {
		os.Setenv("SYNAPSE_URL", "http://synapse")
		os.Setenv("SYNAPSE_TOKEN", "token")
		os.Setenv("USERLI_URL", "http://userli")
		os.Setenv("USERLI_DOMAIN", "example.com")
		os.Setenv("USERLI_TOKEN", "token")

		err := checkEnvironment()

		s.NoError(err)
	})

	s.Run("MissingSynapseURL", func() {
		os.Clearenv()
		err := checkEnvironment()

		s.Error(err)
		s.Contains(err.Error(), "SYNAPSE_URL")
	})

	s.Run("MissingSynapseToken", func() {
		os.Clearenv()

		os.Setenv("SYNAPSE_URL", "http://synapse")
		err := checkEnvironment()

		s.Error(err)
		s.Contains(err.Error(), "SYNAPSE_TOKEN")
	})

	s.Run("MissingUserliURL", func() {
		os.Clearenv()

		os.Setenv("SYNAPSE_URL", "http://synapse")
		os.Setenv("SYNAPSE_TOKEN", "token")
		err := checkEnvironment()

		s.Error(err)
		s.Contains(err.Error(), "USERLI_URL")
	})

	s.Run("MissingUserliDomain", func() {
		os.Clearenv()

		os.Setenv("SYNAPSE_URL", "http://synapse")
		os.Setenv("SYNAPSE_TOKEN", "token")
		os.Setenv("USERLI_URL", "http://userli")
		err := checkEnvironment()

		s.Error(err)
		s.Contains(err.Error(), "USERLI_DOMAIN")
	})

	s.Run("MissingUserliToken", func() {
		os.Clearenv()

		os.Setenv("SYNAPSE_URL", "http://synapse")
		os.Setenv("SYNAPSE_TOKEN", "token")
		os.Setenv("USERLI_URL", "http://userli")
		os.Setenv("USERLI_DOMAIN", "example.com")
		err := checkEnvironment()

		s.Error(err)
		s.Contains(err.Error(), "USERLI_TOKEN")
	})
}

func TestMain(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
