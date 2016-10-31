package main

import (
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type PostgresDatabase struct {
	DB *sql.DB
}

func NewPostgresDatabase(dbHost string, dbPort int, dbUser string, dbPassword string, dbRelName string) (*PostgresDatabase, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%v dbname=%v password=%v port=%v host=%v sslmode=disable",
			dbUser, dbRelName, dbPassword, dbPort, dbHost))
	if err != nil {
		return nil, err
	}
	pg := PostgresDatabase{DB: db}
	return &pg, nil
}

func (db PostgresDatabase) GetUserById(userId string) (User, error) {
	stmt, err := db.DB.Prepare("SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1")
	if err != nil {
		return User{}, nil
	}

	results, err := stmt.Query(userId)
	defer results.Close()
	if err != nil {
		return User{}, err
	}

	hasItem := results.Next()
	if !hasItem {
		return User{}, NotFoundError{"User"}
	}

	var user User
	err = results.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}

	return User{}, nil
}

func (db PostgresDatabase) GetBookmarkByUserId(userId string, itemOffset int, itemLimit int) ([]Bookmark, error) {
	return nil, nil
}

func (db PostgresDatabase) GetBookmarkById(bookmarkId string) (Bookmark, error) {
	stmt, err := db.DB.Prepare("SELECT id, user_id, url, created_at, updated_at WHERE id = $1")
	if err != nil {
		return Bookmark{}, err
	}

	results, err := stmt.Query(bookmarkId)
	defer results.Close()
	if err != nil {
		return Bookmark{}, err
	}

	hasItem := results.Next()
	if !hasItem {
		return Bookmark{}, NotFoundError{"Bookmark"}
	}

	var bookmark Bookmark
	err = results.Scan(&bookmark.Id, &bookmark.User, &bookmark.Url, &bookmark.CreatedAt, &bookmark.UpdatedAt)
	if err != nil {
		return Bookmark{}, err
	}

	tags, err := db.GetTagsForBookmark(bookmarkId)
	if err != nil {
		return Bookmark{}, err
	}
	bookmark.Tags = tags
	return bookmark, nil
}

func (db PostgresDatabase) GetTagsForBookmark(bookmarkId string) ([]string, error) {
	stmt, err := db.DB.Prepare("SELECT tag FROM bookmark_tags WHERE bookmark_id = $1")
	if err != nil {
		return nil, err
	}

	results, err := stmt.Query(bookmarkId)
	defer results.Close()
	if err != nil {
		return nil, err
	}

	var tags []string
	for results.Next() {
		var tag string
		results.Scan(&tag)
		tags = append(tags, tag)
	}

	return tags, nil
}

func (db PostgresDatabase) SaveUser(user User, upsert bool) (User, error) {
	// Check if the user exist
	notfound := false
	_, err := db.GetUserById(uuid.UUID(user.Id).String())

	// Test if the error is not_found
	if err != nil {
		_, notfound = err.(NotFoundError)
	}
	if err != nil && (!notfound || !upsert) { // not found or not upsert or error
		return User{}, err
	}

	// IF user not found, insert
	if notfound {
		stmt, err := db.DB.Prepare("INSERT INTO users(id, username, email, created_at, updated_at) values (gen_random_uuid(), $1, $2, now(), now()) RETURNING id, created_at, updated_at")
		if err != nil {
			return User{}, err
		}
		err = stmt.QueryRow(user.Username, user.Email).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return User{}, err
		}
		return user, nil
	} else {
		stmt, err := db.DB.Prepare("UPDATE TABLE users SET email = $1, updated_at = now() WHERE id = $2 RETURNING updated_at")
		if err != nil {
			return User{}, err
		}
		_, err = stmt.Exec(user.Email, user.Id)
		if err != nil {
			return User{}, err
		}
	}
	return user, nil
}
