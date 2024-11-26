package helpers

import (
	"fmt"
	"log/slog"
)

func WrapErr(s string, e error) error {
	return fmt.Errorf("%s : %w", s, e)
}

func SlErr(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
