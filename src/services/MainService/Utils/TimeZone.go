package Utils

import (
	"fmt"
	"time"
)

func GetTime(location string) (string, error) {
	TimeZone, err := time.LoadLocation(location)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return time.Now().In(TimeZone).Format(time.RFC1123), nil
}
