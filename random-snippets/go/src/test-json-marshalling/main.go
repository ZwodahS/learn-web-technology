package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

/**
This file shows various way to override the json marshalling if you want to change the output of the struct.
For the purpose of this demo, we will try to output time.Time to microsecond timestamp and back.
*/

// -------------------- Method 1 --------------
// http://choly.ca/post/go-json-marshalling/
type MethodOne struct {
	Value      time.Time `json:"timestamp"`
	OtherField string    `json:"other_field"`
}

func (v MethodOne) MarshalJSON() ([]byte, error) {
	type Alias MethodOne
	return json.Marshal(&struct {
		Value int64 `json:"timestamp"`
		*Alias
	}{
		Value: (v.Value.UnixNano() / int64(time.Microsecond)),
		Alias: (*Alias)(&v),
	})
}

func (v *MethodOne) UnmarshalJSON(data []byte) error {
	type Alias MethodOne
	aux := &struct {
		Value int64 `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.Value = time.Unix(0, aux.Value*int64(time.Microsecond))
	return nil
}

// -------------- Method 2-----------------
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

type MethodTwo struct {
	Value      MicrosecondTime `json:"timestamp"`
	OtherField string          `json:"other_field"`
}

func main() {

	data := "{\"timestamp\": 1477932426091642, \"other_field\": \"test\" }"

	var methodOne MethodOne
	var methodTwo MethodTwo

	err := json.Unmarshal([]byte(data), &methodOne)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(data), &methodTwo)
	if err != nil {
		panic(err)
	}

	fmt.Println(methodOne.Value)
	fmt.Println((time.Time)(methodTwo.Value))

	methodOneString, _ := json.Marshal(methodOne)
	methodTwoString, _ := json.Marshal(methodTwo)

	fmt.Println(string(methodOneString))
	fmt.Println(string(methodTwoString))

}
