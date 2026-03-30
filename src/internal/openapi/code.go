package openapi

// Error Codes
const (
	// Generic System Errors
	CodeInternalError    = "INTERNAL_SERVER_ERROR"
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	CodeNotFound         = "RESOURCE_NOT_FOUND"

	// Auth & Signup Errors
	CodeInvalidInput     = "INVALID_INPUT"
	CodeEmailTaken       = "EMAIL_ALREADY_EXISTS"
	CodeWeakPassword     = "PASSWORD_TOO_WEAK"
	CodeUnauthorized     = "UNAUTHORIZED_ACCESS"

	// // Booking Specific Errors
	// CodeRoomUnavailable  = "ROOM_NOT_AVAILABLE"
	// CodeInvalidDates     = "INVALID_BOOKING_DATES"
)
