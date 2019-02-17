package app

import (
	"context"
	"database/sql"
	"fmt"

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

func InitPostgres(ctx context.Context, o DbOptions) (db *sql.DB, err error) {
	// TODO reuse context from init
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

func dbCreateMessage(ctx context.Context, db *sql.DB, m *Message) (err error) {
	query := `INSERT INTO messages(user_id, message_id, is_broadcast, from_admin
			created_at, updated_at, text, reply_to_message, document_id, photo_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err = db.ExecContext(ctx, query, m.UserID, m.MessageID, m.IsBroadcast, m.FromAdmin,
		m.CreatedAt, m.UpdatedAt, m.Text, m.ReplyToMessage, m.DocumentID, m.PhotoID)
	return
}

func dbListUsers(ctx context.Context, db *sql.DB, page int) ([]User, error) {
	// TODO add page support, currently response is unlimited
	// select all users and fetch their last messages if any
	users := make([]User, 0) // 50 is default page size, since there's no real pages, only afterloading
	query := `SELECT DISTINCT users.user_id, users.created_at, users.updated_at,
			users.chat_id, users.email, users.name, users.username,
			message_id, is_broadcast, from_admin, messages.created_at,
			messages.updated_at, text, document_id, photo_id
			FROM users LEFT JOIN messages ON messages.message_id = users.last_message_id`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		u := User{}
		msg := Message{}
		if err := rows.Scan(&u.UserID, &u.CreatedAt, &u.UpdatedAt,
			&u.ChatID, &u.Email, &u.Name, &u.Username,
			&msg.MessageID, &msg.IsBroadcast, &msg.FromAdmin, &msg.CreatedAt,
			&msg.UpdatedAt, &msg.Text, &msg.DocumentID, &msg.PhotoID); err != nil {
			return users, err
		}
		msg.UserID = u.UserID
		u.LastMessage = msg
		users = append(users, u)
	}
	return users, rows.Err()
}

func dbListUserMessages(ctx context.Context, db *sql.DB, userID int, page int) ([]Message, error) {
	messages := make([]Message, 0) // FIXME add actual afterloading of messages, currently whole history is loaded
	// FIXME add ORDER BY clause
	query := `SELECT message_id, is_broadcast, from_admin, created_at,
			updated_at, text, reply_to_message, document_id, photo_id FROM messages WHERE user_id=$1`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return messages, err
	}
	defer rows.Close()
	for rows.Next() {
		msg := Message{}
		if err := rows.Scan(&msg.MessageID, &msg.IsBroadcast, &msg.FromAdmin, &msg.CreatedAt,
			&msg.UpdatedAt, &msg.Text, &msg.ReplyToMessage, &msg.DocumentID, &msg.PhotoID); err != nil {
			return messages, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

// NOTE Use sqlmock database instead
type mockDatabase struct {
}
