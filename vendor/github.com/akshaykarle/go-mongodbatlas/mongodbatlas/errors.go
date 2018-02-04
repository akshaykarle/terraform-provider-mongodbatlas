package mongodbatlas

import (
	"fmt"
)

// APIError represents a MongDB Atlas API Error response
// https://docs.atlas.mongodb.com/api/#errors
type APIError struct {
	Detail    string `json:"detail"`
	Code      int    `json:"error"`
	ErrorCode string `json:"errorCode"`
	Reason    string `json:"reason"`
}

func (e APIError) Error() string {
	if e == (APIError{}) {
		return ""
	}
	return fmt.Sprintf("MongoDB Atlas: %d %v", e.Code, e.Detail)
}

// relevantError returns any non-nil http-related error (creating the request,
// getting the response, decoding) if any. If the decoded apiError is non-nil
// the apiError is returned. Otherwise, no errors occurred, returns nil.
func relevantError(httpError error, apiError APIError) error {
	if httpError != nil {
		return httpError
	}
	if apiError == (APIError{}) {
		return nil
	}
	return apiError
}
