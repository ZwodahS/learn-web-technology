package main

import (
	"fmt"

	"github.com/namsral/flag"

	_ "github.com/lib/pq" // pg
)

func main() {

	var dbPort = flag.Int("dbport", 5432, "database port")
	var dbHost = flag.String("dbstring", "localhost", "database host")
	var dbUser = flag.String("dbuser", "user", "database user")
	var dbPassword = flag.String("dbpassword", "password", "database password")
	var dbName = flag.String("dbname", "url", "database name")

	flag.Parse()

	db, err := NewPostgresDatabase(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)
	PanicIfError(err)

	fmt.Println(db)

}
