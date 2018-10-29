package util

import (
	"fmt"
	. "simple-prepaid-card/constant"
	"time"

)

func CreateUniqueCardNumber(count int) string {
	firstPart := time.Now().Unix()

	secondPart := (count + 1) % INCREMENT_RESET

	return fmt.Sprintf("%010d", firstPart) + fmt.Sprintf("%05d", secondPart)
}

func ParseDateParameter(date string) (*time.Time, error) {
	if date == "" {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02", date)

	return &t, err
}