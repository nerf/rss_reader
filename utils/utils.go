package utils

import (
	"errors"
	"time"
)

var dateFormats = []string{
	time.RFC822,
	time.RFC822Z,
	time.RFC3339,
	time.UnixDate,
	time.RubyDate,
	time.RFC850,
	time.RFC1123Z,
	time.RFC1123,
	time.ANSIC,
	"Mon, January 2 2006 15:04:05 -0700",
	"Mon, Jan 2 2006 15:04:05 -700",
	"Mon, Jan 2 2006 15:04:05 -0700",
}

// Try to detect input layout and parse it
func ParseDate(dateString string) (dateTime time.Time, err error) {
	for _, format := range dateFormats {
		if dateTime, err = time.Parse(format, dateString); err == nil {
			return
		}
	}

	err = errors.New("Unknown date format: " + dateString)
	return
}
