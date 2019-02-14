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

func NewGormDatabase(o DbOptions) (*gorm.DB, error) {
	switch o.DbType {
	case Postgres:
		url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", o.Host, o.Port, o.User, o.DbName, o.Password)
		db, err := gorm.Open("postgres", url)
		return db, err
	default:
		return &gorm.DB{}, fmt.Errorf("This database type is now supported: %v", o.DbType)
	}
}

// NOTE Use sqlmock database instead
// https://github.com/jirfag/go-queryset/blob/master/queryset/queryset_test.go
type mockDatabase struct {

}
