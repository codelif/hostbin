package format

import "time"

func Timestamp(value string) string {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}

	return parsed.Local().Format("2006-01-02 15:04")
}
