package main

import (
	"database/sql"
	"testing"

	"github.com/namsral/flag"
	uuid "github.com/satori/go.uuid"
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
	var idString string
	err = db.DB.QueryRow("SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1", uuid.UUID(user.Id).String()).Scan(
		&idString, &createdUser.Username, &createdUser.Email, &createdUser.CreatedAt, &createdUser.UpdatedAt)
	PanicIfError(err)

	u, err := uuid.FromString(idString)
	PanicIfError(err)

	createdUser.Id = ShortUUID(u)
	if createdUser.Id != user.Id {
		t.Error("Created User id and DB user id is different")
	}

	Teardown(db.DB, t)
}
