package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SynapseUserResponse represents a response from the Synapse API
type SynapseUserResponse struct {
	Users []SynapseUser `json:"users"`
}

// SynapseUser represents a user in Synapse
type SynapseUser struct {
	Name      string `json:"name"`
	CreatedAt int64  `json:"creation_ts"`
	LastSeen  int64  `json:"last_seen_ts"`
}

// SynapseClient is a client for the Synapse API
type SynapseClient struct {
	url   string
	token string

	client *http.Client
}

// NewSynapseClient creates a new Synapse client
func NewSynapseClient(url, token string) *SynapseClient {
	return &SynapseClient{
		url:   url,
		token: token,

		client: &http.Client{},
	}
}

// FetchUsers fetches all users from Synapse
func (s *SynapseClient) FetchUsers() ([]SynapseUser, error) {
	url := fmt.Sprintf("%s/_synapse/admin/v3/users?deactivated=false&guests=false&limit=%d", s.url, 99999999)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var body SynapseUserResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Users, nil
}

// doRequest performs an HTTP request with the client
func (s *SynapseClient) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Userli-Synapse-User-Retention")

	return s.client.Do(req)
}
