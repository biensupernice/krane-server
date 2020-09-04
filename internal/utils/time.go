package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

// UTCDateString : Get the current date time in RFC3339 format
func UTCDateString() string {
	t := time.Now().Local()
	return t.Format(time.RFC3339)
}

func MakeIdentifier() string {
	b := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}