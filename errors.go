package clubhouse

import "github.com/pkg/errors"

// Custom errors.
var (
	ErrSignatureMismatch = errors.New("Signature mismatch")
)
