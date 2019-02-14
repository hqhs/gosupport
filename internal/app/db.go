package app

// Database implements db access methods.
type Database interface {

}

// MockDatabase implements Database interface and stores data
// in-memory for testing/prototyping
type MockDatabase struct {

}

// NewMockDatabase initializes database connection
func NewMockDatabase() Database {
	return &MockDatabase{}
}
