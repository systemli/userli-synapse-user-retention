package main

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type SynapseTestSuite struct {
	suite.Suite
}

func (s *SynapseTestSuite) SetupTest() {
	defer gock.Off()
}

func (s *SynapseTestSuite) TestFetchUsers() {
	s.Run("Success", func() {
		gock.New("http://synapse").
			Get("/_synapse/admin/v3/users").
			MatchHeader("Authorization", "Bearer token").
			MatchHeader("Content-Type", "application/json").
			MatchHeader("User-Agent", "Userli-Synapse-User-Retention").
			Reply(200).
			JSON(map[string]interface{}{
				"users": []map[string]interface{}{
					{
						"name":         "test",
						"creation_ts":  123456,
						"last_seen_ts": 123456,
					},
				},
			})

		client := NewSynapseClient("http://synapse", "token")
		users, err := client.FetchUsers()

		s.True(gock.IsDone())
		s.NoError(err)
		s.Len(users, 1)
		s.Equal("test", users[0].Name)
	})

	s.Run("Synapse Error", func() {
		gock.New("http://synapse").
			Get("/_synapse/admin/v3/users").
			Reply(500)

		client := NewSynapseClient("http://synapse", "token")
		users, err := client.FetchUsers()

		s.True(gock.IsDone())
		s.Error(err)
		s.Nil(users)
	})

	s.Run("Network Error", func() {
		client := NewSynapseClient("http://synapse", "token")
		users, err := client.FetchUsers()

		s.Error(err)
		s.Nil(users)
	})
}

func TestSynapseClient(t *testing.T) {
	suite.Run(t, new(SynapseTestSuite))
}
