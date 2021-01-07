package clubhouse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// CHClient defines available methods.
type CHClient interface {
	CreateEpic(name, description string) (Epic, error)
}

const apiURL = "https://api.clubhouse.io/api/v3"

// Options are the settings needed when creating a new Client.
type Options struct {
	Token string
}

// Client handles interaction with the Clubhouse API.
type Client struct {
	url   string
	token string
}

// NewClient creates and returns a new Clubhouse Client.
func NewClient(options Options) CHClient {
	client := Client{
		token: options.Token,
		url:   apiURL,
	}

	return client
}

// CreateEpic creates an Epic on Clubhouse.
func (c Client) CreateEpic(name, description string) (Epic, error) {
	epic := Epic{
		Name:        name,
		Description: description,
	}

	body, err := json.Marshal(epic)
	if err != nil {
		return epic, errors.Wrap(err, "Could not create Epic body")
	}

	httpClient := &http.Client{}

	url := c.url + "/epics"

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	request.Header.Add("Clubhouse-Token", c.token)
	request.Header.Add("Content-Type", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return epic, errors.Wrap(err, "Could not send request to the Clubhouse API")
	}

	if response.StatusCode != http.StatusCreated {
		return epic, errors.New(fmt.Sprint("Clubhouse returned status", response.StatusCode))
	}

	return epic, nil
}
