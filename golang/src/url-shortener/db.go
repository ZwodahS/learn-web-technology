package main

import (
	"database/sql"
	"fmt"
)

type DatabaseWrapper struct {
	DB *sql.DB
}

func (self *DatabaseWrapper) UpdateUrl(url Url) (*Url, error) {
	existingUrl, err := self.GetURLById(url.Id)
	if err != nil {
		return nil, err
	}
	fmt.Println(err)

	if existingUrl == nil {
		stmt, err := self.DB.Prepare("INSERT INTO url(id, url) values (?, ?)")
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(url.Id, url.Url)
		defer stmt.Close()
		if err != nil {
			return nil, err
		}
	} else {
		stmt, err := self.DB.Prepare("UPDATE url SET url = ? WHERE id = ?")
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(url.Url, url.Id)
		defer stmt.Close()
		if err != nil {
			return nil, err
		}
	}
	return &url, nil
}

func (self *DatabaseWrapper) GetURLById(urlId string) (*Url, error) {
	stmt, err := self.DB.Prepare("SELECT id, url FROM url WHERE id=?")
	if err != nil {
		panic(err)
		//return nil, err
	}

	// url := Url{Id: "", Url: ""}
	var url Url
	results, err := stmt.Query(urlId)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	hasItem := results.Next()
	if !hasItem {
		return nil, nil
	}
	results.Scan(&url.Id, &url.Url)
	return &url, nil
}
