package main

import (
	"fmt"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

// A Microsecond time for us to marshal json
type MicrosecondTime time.Time

func (t *MicrosecondTime) UnmarshalJSON(data []byte) error {
	microseconds, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = MicrosecondTime(time.Unix(0, microseconds*int64(time.Microsecond)))
	return nil
}

func (t MicrosecondTime) MarshalJSON() ([]byte, error) {
	_t := time.Time(t)
	return []byte(fmt.Sprintf("%v", (_t.UnixNano() / int64(time.Microsecond)))), nil
}

type ShortUUID uuid.UUID

type User struct {
	Id        ShortUUID       `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	CreatedAt MicrosecondTime `json:"created_at"`
	UpdatedAt MicrosecondTime `json:"updated_at"`
}

type Bookmark struct {
	Id        ShortUUID `json:"id"`
	User      string
	Url       string          `json:"url"`
	Tags      []string        `json:"tags"`
	CreatedAt MicrosecondTime `json:"created_at"`
	UpdatedAt MicrosecondTime `json:"updated_at"`
}
