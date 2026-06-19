package dnse

import "fmt"

// APIError is returned when the DNSE server responds with a non-2xx status code.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("dnse: API error %d: %s", e.StatusCode, e.Body)
}
