package main

import (
	"database/sql"
	"flag"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/renstrom/shortuuid"
	uuid "github.com/satori/go.uuid"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type Book struct {
	Id        uuid.UUID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {

	var dbPort = flag.Int("dbport", 5432, "database port")
	var dbHost = flag.String("dbstring", "localhost", "database host")
	var dbUser = flag.String("dbuser", "user", "database user")
	var dbPassword = flag.String("dbpassword", "password", "database password")
	var dbName = flag.String("dbname", "test_uuid", "database name")
	flag.Parse()

	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%v dbname=%v password=%v port=%v host=%v sslmode=disable",
		*dbUser, *dbName, *dbPassword, *dbPort, *dbHost))
	checkError(err)
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM books")
	_, err = stmt.Exec()
	checkError(err)
	stmt, err = db.Prepare("INSERT INTO books(id, title, created_at, updated_at) values ( gen_random_uuid(), $1, now(), now())")
	checkError(err)

	_, err = stmt.Exec("Introduction to Go")
	checkError(err)
	_, err = stmt.Exec("Advanced Go Programming")
	checkError(err)
	defer stmt.Close()

	stmt, err = db.Prepare("SELECT id, title, created_at, updated_at FROM books")
	checkError(err)
	result, err := stmt.Query()
	checkError(err)
	for result.Next() {
		var book Book
		var idstring string
		result.Scan(&idstring, &book.Title, &book.CreatedAt, &book.UpdatedAt)
		u2, err := uuid.FromString(idstring)
		checkError(err)
		book.Id = u2
		short := shortuuid.DefaultEncoder.Encode(book.Id)
		longagain, _ := shortuuid.DefaultEncoder.Decode(short)
		fmt.Printf("UUID: %v, Shortuuid: %v, Long again %v, Title: %v\n", book.Id, short, longagain, book.Title)
		fmt.Println(book.CreatedAt.UnixNano(), "                  ", book.UpdatedAt.UnixNano())
	}
}
