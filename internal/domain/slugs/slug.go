package slugs

import (
	"errors"
	"regexp"
)

var (
	ErrInvalid  = errors.New("invalid slug")
	ErrReserved = errors.New("reserved slug")

	slugPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)
)

func Validate(value string, reserved map[string]struct{}) error {
	if !slugPattern.MatchString(value) {
		return ErrInvalid
	}

	if _, ok := reserved[value]; ok {
		return ErrReserved
	}

	return nil
}
