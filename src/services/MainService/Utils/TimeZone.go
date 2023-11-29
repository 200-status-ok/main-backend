package Utils

import (
	"fmt"
	"time"
)

func GetTime(location string) (int64, error) {
	TimeZone, err := time.LoadLocation(location)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return time.Now().In(TimeZone).Unix(), nil
}
