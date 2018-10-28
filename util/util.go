package util

import (
	"fmt"
	. "prepaid-card/constant"
	"time"

)

func CreateUniqueCardNumber(count int) string {
	firstPart := time.Now().Unix()

	secondPart := (count + 1) % INCREMENT_RESET

	return fmt.Sprintf("%010d", firstPart) + fmt.Sprintf("%05d", secondPart)
}