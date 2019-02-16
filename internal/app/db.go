package app

import (
	"fmt"
	"context"
	"database/sql"

	_ "github.com/lib/pq" // postgres driver
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
}


func InitPostgres(o DbOptions) (db *sql.DB, err error) {
	// TODO reuse context from init
	ctx := context.Background()
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		o.Host, o.Port, o.User, o.DbName, o.Password)
	fmt.Println("url: ", url)
	// TODO init statements
	db, err = sql.Open("postgres", url)
	if err := db.PingContext(ctx); err != nil {
		return db, err
	}
	return
}

func dbCreateAdmin(ctx context.Context, db *sql.DB, a *Admin) (err error) {
	query := `INSERT INTO admins(created_at, updated_at,
			email, name, hashed_password, is_superuser,
			is_active, email_confirmed, auth_token, password_reset_token)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err = db.ExecContext(ctx, query, a.CreatedAt, a.UpdatedAt,
		a.Email, a.Name, a.HashedPassword, a.IsSuperUser,
		a.IsActive, a.EmailConfirmed, a.AuthToken, a.PasswordResetToken)
	return
}



// NOTE Use sqlmock database instead
type mockDatabase struct {

}
