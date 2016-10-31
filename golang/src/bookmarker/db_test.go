package main

import (
	"database/sql"
	"testing"

	"github.com/namsral/flag"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func NewDBFromEnvVar(t *testing.T) *PostgresDatabase {
	var dbPort = flag.Int("dbport", 5432, "database port")
	var dbHost = flag.String("dbstring", "localhost", "database host")
	var dbUser = flag.String("dbuser", "user", "database user")
	var dbPassword = flag.String("dbpassword", "password", "database password")
	var dbName = flag.String("dbname", "url", "database name")

	flag.Parse()
	db, err := NewPostgresDatabase(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)
	PanicIfError(err)

	return db
}

func Setup(db *sql.DB, t *testing.T) {
	_, err := db.Exec("DELETE FROM users")
	PanicIfError(err)

	_, err = db.Exec("DELETE FROM bookmarks")
	PanicIfError(err)

	_, err = db.Exec("DELETE FROM bookmark_tags")
	PanicIfError(err)
}

func Teardown(db *sql.DB, t *testing.T) {
	_, err := db.Exec("DELETE FROM users")
	PanicIfError(err)

	_, err = db.Exec("DELETE FROM bookmarks")
	PanicIfError(err)

	_, err = db.Exec("DELETE FROM bookmark_tags")
	PanicIfError(err)
}

func TestCreateUser(t *testing.T) {
	db := NewDBFromEnvVar(t)
	Setup(db.DB, t)

	user := User{Username: "testuser", Email: "test@hello.com"}

	user, err := db.SaveUser(user, true)
	PanicIfError(err)

	var createdUser User
	stmt, err := db.DB.Prepare("SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1")
	result, err := stmt.Query(uuid.UUID(user.Id).String())
	assert.True(t, result.Next())
	result.Scan(&createdUser.Id, &createdUser.Username, &createdUser.Email, &createdUser.CreatedAt, &createdUser.UpdatedAt)
	PanicIfError(err)

	assert.Equal(t, createdUser, user)

	Teardown(db.DB, t)
}
