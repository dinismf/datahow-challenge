package core

import "fmt"

// Pre-defined ServiceErrors shared across the service layer.
// Use WithReason to attach internal context without changing the client message:
//
//	return core.ErrSvcNotFound.WithReason(fmt.Errorf("flagID=%s", id))
var (
	ErrSvcNotFound     = &ServiceError{Code: "4000", Message: "resource not found"}
	ErrSvcInvalidInput = &ServiceError{Code: "4001", Message: "invalid input"}
	ErrSvcConflict     = &ServiceError{Code: "4002", Message: "resource already exists"}
	ErrSvcInternal     = &ServiceError{Code: "5000", Message: "internal error"}
)

// ServiceError is returned by the service layer.
// Message is safe to send to clients; Reason holds internal context for logging only.
type ServiceError struct {
	Code    string
	Message string
	Reason  error // never serialised; use LogError() to access
}

// Error returns the client-safe message only. Reason is never included.
func (e *ServiceError) Error() string { return e.Message }

// LogError returns the full string including the Reason chain.
// Use only for server-side logging, never in responses.
func (e *ServiceError) LogError() string {
	if e.Reason != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Reason)
	}
	return e.Message
}

// NewServiceError constructs a ServiceError for business-logic violations
// where no underlying infrastructure error is being wrapped.
func NewServiceError(code string, message string) *ServiceError {
	return &ServiceError{Code: code, Message: message}
}

// WithReason returns a copy of the error with an internal reason attached. This is a useful context for logging purposes.
// Returns a new value so catalog sentinels are never mutated.
func (e *ServiceError) WithReason(reason error) *ServiceError {
	return &ServiceError{Code: e.Code, Message: e.Message, Reason: reason}
}
