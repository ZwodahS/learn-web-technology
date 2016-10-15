package main

import (
	"fmt"
	"strconv"
	"time"
)

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

type Url struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
