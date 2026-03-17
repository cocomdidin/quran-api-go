package validator

// Validate checks that id is correct.

import (
	"math"
	"quran-api-go/internal/domain"
	"strconv"
)

// Returns domain.ErrInvalidIDParam for any other value.
//
// Usage:
//
//	lang, err := validator.ValidateIDParam(c.Query("id"))
//	if err != nil {
//	    response.BadRequest(c, "invalid param")
//	    return
//	}
func ValidateIDParam(id string) (string, error) {
	convertedId, err := strconv.Atoi(id)
	if err != nil {
		return "", domain.ErrInvalidIDParam
	}

	if math.IsNaN(float64(convertedId)) {
		return "", domain.ErrInvalidIDParam
	}

	if convertedId < 1 {
		return "", domain.ErrInvalidIDParam
	}

	return id, nil
}

// Returns domain.ErrInvalidRangeParam for any other value.
//
// Usage:
//
//	err := validator.ValidateRangeParam(c.Query("from"), c.Query("to"))
//	if err != nil {
//	    response.BadRequest(c, "invalid param")
//	    return
//	}
func ValidateRangeParam(from, to string) error {
	convertedFrom, err := strconv.Atoi(from)
	if err != nil {
		return domain.ErrInvalidRangeParam
	}

	convertedTo, err := strconv.Atoi(to)
	if err != nil {
		return domain.ErrInvalidRangeParam
	}

	// if one of it contain NaN value
	if math.IsNaN(float64(convertedFrom)) || math.IsNaN(float64(convertedTo)) {
		return domain.ErrInvalidRangeParam
	}

	if convertedFrom < 1 || convertedTo < 1 {
		return domain.ErrInvalidRangeParam
	}

	if convertedFrom > convertedTo {
		return domain.ErrInvalidRangeParam
	}

	return nil
}
