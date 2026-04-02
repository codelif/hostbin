package format

import "github.com/dustin/go-humanize"

func Bytes(size int64) string {
	return humanize.Bytes(uint64(size))
}
