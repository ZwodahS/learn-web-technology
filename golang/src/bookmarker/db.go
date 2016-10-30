package main

import "fmt"

type NotFoundError struct {
	ObjectType string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%v not found", e.ObjectType)
}

type Database interface {
	GetUserById(userId string) (User, error)
	GetBookmarkByUserId(userId string, itemOffset int, itemLimit int) ([]Bookmark, error)
	GetBookmarkById(bookmarkId string) (Bookmark, error)
	SaveUser(user User, upsert bool) error
}
