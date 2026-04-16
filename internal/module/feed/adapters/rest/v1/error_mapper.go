// Package v1 provides REST API v1 handlers for the feed module.
package v1

import (
	"errors"

	"real-time-activity-feed/internal/module/feed/domain"
	"real-time-activity-feed/internal/shared/response"
	"real-time-activity-feed/internal/shared/validator"
)

// toAPIError converts feed errors to APIError (internal helper)
func toAPIError(err error) *response.APIError {
	if err == nil {
		return nil
	}

	// Check for validation errors
	if validator.IsValidationError(err) {
		var validationErr *validator.ValidationError
		if errors.As(err, &validationErr) {
			return response.NewValidationError(validationErr.Message)
		}
	}

	var invalidEventTypeErr *domain.InvalidEventTypeError
	if errors.As(err, &invalidEventTypeErr) {
		return response.NewValidationError(invalidEventTypeErr.Error())
	}

	// If it's already an APIError, return it as-is
	if apiErr, ok := err.(*response.APIError); ok {
		return apiErr
	}

	// Default to internal error for wrapped errors
	return response.NewInternalError("An error occurred")
}
