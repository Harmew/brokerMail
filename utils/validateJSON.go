package utils

import (
	"errors"
	"github.com/Harmew/brokerMail/models"
)

func ValidateJSON(interfaceBody models.SendGridInternal) error {
	if interfaceBody.Subject == "" {
		return errors.New("Field 'subject' is required")
	}

	if len(interfaceBody.Recipients) == 0 {
		return errors.New("Field 'recipients' is required")
	}

	if interfaceBody.Content == "" {
		return errors.New("Field 'content' is required")
	}

	return nil
}
