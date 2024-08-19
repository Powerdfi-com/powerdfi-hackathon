package helpers

import (
	"net/mail"

	"github.com/google/uuid"
)

func IsValidUuid(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}

func IsStatusOk(statusCode int) bool {
	statusOK := statusCode >= 200 && statusCode < 300
	return statusOK
}

func IsValidEmail(value string) bool {
	_, err := mail.ParseAddress(value)
	return err == nil
}
