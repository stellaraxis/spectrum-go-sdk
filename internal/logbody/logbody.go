package logbody

import "unicode/utf8"

const MaxBytes = 8192

type Result struct {
	Message       string
	Truncated     bool
	OriginalBytes int
	MaxBytes      int
}

func Normalize(message string) Result {
	originalBytes := len(message)
	result := Result{
		Message:       message,
		OriginalBytes: originalBytes,
		MaxBytes:      MaxBytes,
	}

	if originalBytes <= MaxBytes {
		return result
	}

	truncated := message[:MaxBytes]
	for len(truncated) > 0 && !utf8.ValidString(truncated) {
		truncated = truncated[:len(truncated)-1]
	}

	result.Message = truncated
	result.Truncated = true
	return result
}
