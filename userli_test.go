package main

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type UserliTestSuite struct {
	suite.Suite
}

func (s *UserliTestSuite) SetupTest() {
	defer gock.Off()
}

func (s *UserliTestSuite) TestFetchDeletedUsers() {
	s.Run("Success", func() {
		gock.New("http://userli").
			Get("/api/retention/test/users").
			MatchHeader("Authorization", "Bearer token").
			MatchHeader("Content-Type", "application/json").
			MatchHeader("User-Agent", "Userli-Synapse-User-Retention").
			Reply(200).
			JSON([]string{"test"})

		client := NewUserliClient("http://userli", "test", "token")
		users, err := client.FetchDeletedUsers()

		s.True(gock.IsDone())
		s.NoError(err)
		s.Len(users, 1)
		s.Equal("test", users[0])
	})

	s.Run("Userli Error", func() {
		gock.New("http://userli").
			Get("/api/retention/test/users").
			Reply(500)

		client := NewUserliClient("http://userli", "test", "token")
		users, err := client.FetchDeletedUsers()

		s.True(gock.IsDone())
		s.Error(err)
		s.Nil(users)
	})

	s.Run("Network Error", func() {
		client := NewUserliClient("http://userli", "test", "token")
		users, err := client.FetchDeletedUsers()

		s.Error(err)
		s.Nil(users)
	})
}

func (s *UserliTestSuite) TestTouchUser() {
	s.Run("Success", func() {
		gock.New("http://userli").
			Put("/api/retention/test/touch").
			MatchHeader("Authorization", "Bearer token").
			MatchHeader("Content-Type", "application/json").
			MatchHeader("User-Agent", "Userli-Synapse-User-Retention").
			Reply(200)

		client := NewUserliClient("http://userli", "test", "token")
		err := client.TouchUser("test", 123456)

		s.True(gock.IsDone())
		s.NoError(err)
	})

	s.Run("Userli Error", func() {
		gock.New("http://userli").
			Put("/api/retention/test/touch").
			Reply(500)

		client := NewUserliClient("http://userli", "test", "token")
		err := client.TouchUser("test", 123456)

		s.True(gock.IsDone())
		s.Error(err)
	})

	s.Run("Network Error", func() {
		client := NewUserliClient("http://userli", "test", "token")
		err := client.TouchUser("test", 123456)

		s.Error(err)
	})
}

func TestUserliClient(t *testing.T) {
	suite.Run(t, new(UserliTestSuite))
}
