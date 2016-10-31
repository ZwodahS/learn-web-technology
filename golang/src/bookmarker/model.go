package main

import (
	"encoding/json"
	"time"

	"github.com/renstrom/shortuuid"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// http://choly.ca/post/go-json-marshalling/
func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Id        string `json:"id"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
		*Alias
	}{
		Id:        shortuuid.DefaultEncoder.Encode(u.Id),
		CreatedAt: (u.CreatedAt.UnixNano() / int64(time.Microsecond)),
		UpdatedAt: (u.UpdatedAt.UnixNano() / int64(time.Microsecond)),
		Alias:     (*Alias)(&u),
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		Id        string `json:"id"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	u.Id, err = shortuuid.DefaultEncoder.Decode(aux.Id)
	if err != nil {
		return err
	}

	u.CreatedAt = time.Unix(0, aux.CreatedAt*int64(time.Microsecond))
	u.UpdatedAt = time.Unix(0, aux.UpdatedAt*int64(time.Microsecond))

	return nil
}

type Bookmark struct {
	Id        uuid.UUID `json:"id"`
	User      string
	Url       string    `json:"url"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
