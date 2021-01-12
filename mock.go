package clubhouse

import (
	"io"

	"github.com/pkg/errors"
)

// ErrTest is returned for testing method errors.
var ErrTest = errors.New("Test error")

// MockClient is a mock Client to use for testing.
type MockClient struct {
	Epic              Epic
	Webhook           Webhook
	CreateEpicError   bool
	ParseWebhookError bool
}

// CreateEpic mock creates an Epic on Clubhouse.
func (c MockClient) CreateEpic(name, description string) (Epic, error) {
	if c.CreateEpicError {
		return Epic{}, ErrTest
	}

	return c.Epic, nil
}

// ParseWebhook mock parses a Webhook's body and returns a Webhook struct.
func (c MockClient) ParseWebhook(body io.ReadCloser) (Webhook, error) {
	if c.ParseWebhookError {
		return Webhook{}, errors.Wrap(ErrTest, "Could not parse Webhook body")
	}

	return c.Webhook, nil
}
