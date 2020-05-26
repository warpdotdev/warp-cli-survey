package store

import (
	"time"

	"github.com/zachlloyd/denver-survey-client/shell"
)

// Storer is an interface for recording responses
type Storer interface {
	Write(response Response)
}
