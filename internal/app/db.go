package app

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/jinzhu/gorm.v1"
)

// DbType describes supported databases
type DbType int

const (
	Postgres DbType = iota
	Sqllite
)

type DbOptions struct {
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
	DbType   DbType
}

// Database implements db access methods.
type Database interface {
}

type GormDatabase struct {
	db *gorm.DB
}

func NewGormDatabase(o DbOptions) (Database, error) {
	switch o.DbType {
	case Postgres:
		url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", o.Host, o.Port, o.User, o.DbName, o.Password)
		db, err := gorm.Open("postgres", url)
		return &GormDatabase{db}, err
	default:
		return &MockDatabase{}, fmt.Errorf("This database type is now supported: %v", o.DbType)
	}
}

// MockDatabase implements Database interface and stores data in-memory for testing/prototyping
// NOTE: maybe use bbolt db? No requirements, simple one file storage, consistent and speed is ok
// I'm going to implement db support with gorm first, to get some sql experience in go, and
// skip testing at all. So MockDatabase is useless right now
type MockDatabase struct {
}

// NewMockDatabase initializes database connection
func NewMockDatabase() Database {
	return &MockDatabase{}
}
