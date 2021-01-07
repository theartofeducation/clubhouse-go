package clubhouse

import "github.com/pkg/errors"

// ErrTest is returned for testing method errors.
var ErrTest = errors.New("Test error")

// MockClient is a mock Client to use for testing.
type MockClient struct {
	Epic            Epic
	CreateEpicError bool
}

// CreateEpic mock creates an Epic on Clubhouse.
func (c MockClient) CreateEpic(name, description string) (Epic, error) {
	if c.CreateEpicError {
		return Epic{}, ErrTest
	}

	return c.Epic, nil
}
