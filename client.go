package clubhouse

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// CHClient defines available methods.
type CHClient interface {
	CreateEpic(name, description string) (Epic, error)
	ParseWebhook(body io.ReadCloser) (Webhook, error)
	VerifySignature(signature string, body []byte) error
}

const apiURL = "https://api.clubhouse.io/api/v3"

// Options are the settings needed when creating a new Client.
type Options struct {
	Token         string
	WebhookSecret string
}

// Client handles interaction with the Clubhouse API.
type Client struct {
	url           string
	token         string
	webhookSecret string
}

// NewClient creates and returns a new Clubhouse Client.
func NewClient(options Options) CHClient {
	client := Client{
		token:         options.Token,
		url:           apiURL,
		webhookSecret: options.WebhookSecret,
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

// ParseWebhook parses a Webhook's body and returns a Webhook struct.
func (c Client) ParseWebhook(body io.ReadCloser) (Webhook, error) {
	defer body.Close()

	var webhook Webhook

	if err := json.NewDecoder(body).Decode(&webhook); err != nil {
		return webhook, errors.Wrap(err, "Could not parse Webhook body")
	}

	return webhook, nil
}

// VerifySignature validates a Webhook's signature.
func (c Client) VerifySignature(signature string, body []byte) error {
	secret := []byte("testsecret")

	hash := hmac.New(sha256.New, secret)
	hash.Write(body)
	generatedSignature := hex.EncodeToString(hash.Sum(nil))

	if signature == generatedSignature {
		return nil
	}

	return errors.New("Signature mismatch")
}
