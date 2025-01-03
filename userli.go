package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// UserliClient is a client for the Userli API
type UserliClient struct {
	url    string
	domain string
	token  string

	client *http.Client
}

// NewUserliClient creates a new Userli client
func NewUserliClient(url, domain, token string) *UserliClient {
	return &UserliClient{
		url:    url,
		domain: domain,
		token:  token,

		client: &http.Client{},
	}
}

// FetchDeletedUsers fetches the list of deleted users from Userli
func (c *UserliClient) FetchDeletedUsers() ([]string, error) {
	url := fmt.Sprintf("%s/api/retention/%s/users", c.url, c.domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var emails []string

	if err := json.NewDecoder(res.Body).Decode(&emails); err != nil {
		return nil, err
	}

	return emails, nil
}

// TouchUser touches a user in Userli if the last seen is newer
func (c *UserliClient) TouchUser(email string, timestamp int64) error {
	url := fmt.Sprintf("%s/api/retention/%s/touch", c.url, email)
	body := fmt.Sprintf(`{"timestamp": %d}`, timestamp/1000)
	req, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		return err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("touch user failed with status code %d", res.StatusCode)
	}

	return nil
}

// doRequest performs a request with the correct headers
func (c *UserliClient) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Userli-Synapse-User-Retention")

	return c.client.Do(req)
}
