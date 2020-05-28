package store

// Storer is an interface for recording responses
type Storer interface {
	Write(response Response)
}
