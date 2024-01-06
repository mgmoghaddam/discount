package serr

import (
	"database/sql"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrInternal              ErrorCode = "INTERNAL"
	ErrInvalidGiftID         ErrorCode = "INVALID_GIFT_ID"
	ErrInvalidGiftCode       ErrorCode = "INVALID_GIFT_CODE"
	ErrInvalidDiscountID     ErrorCode = "INVALID_DISCOUNT_ID"
	ErrInvalidDiscountCode   ErrorCode = "INVALID_DISCOUNT_CODE"
	ErrPermission            ErrorCode = "PERMISSION"
	ErrGiftUsageLimitReached ErrorCode = "GIFT_USAGE_LIMIT_REACHED"
)

type ServiceError struct {
	Method    string
	Cause     error
	Message   string
	ErrorCode ErrorCode
	Code      int
}

func (e ServiceError) Error() string {
	return fmt.Sprintf(
		"%s (%d) - %s: %s",
		e.Method, e.Code, e.Message, e.Cause,
	)
}

func ValidationErr(method, message string, code ErrorCode) error {
	return &ServiceError{
		Method:    method,
		Message:   message,
		Code:      http.StatusBadRequest,
		ErrorCode: code,
	}
}

func DBError(method, repo string, cause error) error {
	err := &ServiceError{
		Method: fmt.Sprintf("%s.%s", repo, method),
		Cause:  cause,
	}
	switch cause {
	case sql.ErrNoRows:
		err.Code = http.StatusNotFound
		err.Message = fmt.Sprintf("%s not found", repo)
	default:
		err.Code = http.StatusInternalServerError
		err.Message = fmt.Sprintf("could not perform action on %s", repo)
	}
	return err
}
